package session

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"math/big"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/player"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"
)

type (
	unansweredQuestion struct {
		ID    uuid.UUID
		IsNew bool
	}

	AnswerChecker interface {
		Check(question *model.Question, answers []string) (*contracts.AcceptAnswersOut, error)
	}

	Usecase struct {
		sessions  session.Repository
		games     game.Repository
		questions question.Repository
		players   player.Repository
		template  transactional.Template

		checkers map[model.QuestionType]AnswerChecker
	}
)

func NewUsecase(
	sessions session.Repository,
	games game.Repository,
	questions question.Repository,
	players player.Repository,
	template transactional.Template,
	checkers map[model.QuestionType]AnswerChecker,
) contracts.SessionUsecase {
	return &Usecase{
		sessions:  sessions,
		games:     games,
		questions: questions,
		players:   players,
		template:  template,
		checkers:  checkers,
	}
}

func (u *Usecase) Start(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		if _, err := u.getActiveGame(ctx, tx, gameID); err != nil {
			return err
		}

		specificPlayers, err := u.players.GetByIDs(ctx, []uuid.UUID{playerID})
		if err != nil {
			return err
		}
		if len(specificPlayers) == 0 {
			err = u.players.Insert(ctx, tx, &model.Player{
				ID:   playerID,
				Name: "unknown",
			})
			if err != nil {
				return err
			}
		}

		return u.sessions.Insert(
			ctx,
			tx,
			&model.Session{
				PlayerID: playerID,
				GameID:   gameID,
				Status:   model.SessionStatusStarted,
			},
		)
	})
}

func (u *Usecase) Finish(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		if _, err := u.getActiveGame(ctx, tx, gameID); err != nil {
			return err
		}

		specificPlayerGame, err := u.sessions.GetBySpecWithTx(ctx, tx, &session.Spec{
			PlayerID: playerID,
			GameID:   gameID,
		})
		if err != nil {
			return err
		}

		specificPlayerGame.Status = model.SessionStatusFinished
		return u.sessions.Update(ctx, tx, specificPlayerGame)
	})
}

func (u *Usecase) GetStatistics(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*model.SessionStatistics, error) {
	var result *model.SessionStatistics
	return result, u.template.Execute(ctx, func(tx transactional.Tx) error {
		sessionItems, err := u.sessions.GetSessionBySpecWithTx(
			ctx,
			tx,
			&session.ItemSpec{
				PlayerID: playerID,
				GameID:   gameID,
			},
		)
		if err != nil {
			return err
		}
		if len(sessionItems) == 0 {
			return contracts.ErrSessionNotFound
		}

		totalQuestions := int64(len(sessionItems))
		correctAnswers := int64(0)
		for _, item := range sessionItems {
			if item.IsCorrect == nil || !*item.IsCorrect {
				continue
			}

			correctAnswers++
		}

		result = &model.SessionStatistics{
			QuestionsCount:      totalQuestions,
			CorrectAnswersCount: correctAnswers,
		}
		return nil
	})
}

func (u *Usecase) GetSessions(ctx context.Context, gameID uuid.UUID) ([]model.SessionExtended, error) {
	sessions, err := u.sessions.GetSessionsByGameID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (u *Usecase) getActiveGame(ctx context.Context, tx transactional.Tx, gameID uuid.UUID) (*model.Game, error) {
	specificGame, err := u.games.GetWithTx(ctx, tx, gameID)
	if err != nil {
		return nil, err
	}
	if specificGame.Status != model.GameStatusStarted {
		return nil, errors.New("game isn't started")
	}

	return specificGame, nil
}

func (u *Usecase) getSession(ctx context.Context, tx transactional.Tx, playerID uuid.UUID, gameID uuid.UUID) (*model.Session, error) {
	specificSession, err := u.sessions.GetBySpecWithTx(ctx, tx, &session.Spec{
		PlayerID: playerID,
		GameID:   gameID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		err := u.Start(ctx, gameID, playerID)
		if err != nil {
			return nil, err
		}

		return u.sessions.GetBySpecWithTx(ctx, tx, &session.Spec{
			PlayerID: playerID,
			GameID:   gameID,
		})

	}
	if err != nil {
		return nil, err
	}

	return specificSession, nil
}

func getUnansweredQuestionID(game *model.Game, gameQuestions []uuid.UUID, sessionItems []model.SessionItem) (*unansweredQuestion, error) {
	unansweredSessionItems := slices.Filter(sessionItems, func(item model.SessionItem) bool {
		return item.AnsweredAt == nil
	})
	if len(unansweredSessionItems) > 0 {
		return &unansweredQuestion{
			ID:    unansweredSessionItems[0].QuestionID,
			IsNew: false,
		}, nil
	}

	answeredSessionItemsMap := make(map[uuid.UUID]struct{}, len(sessionItems))
	for _, item := range sessionItems {
		answeredSessionItemsMap[item.QuestionID] = struct{}{}
	}

	unansweredQuestionIDs := slices.Filter(gameQuestions, func(questionID uuid.UUID) bool {
		_, ok := answeredSessionItemsMap[questionID]
		return !ok
	})
	if len(unansweredQuestionIDs) == 0 {
		return nil, contracts.ErrQuestionQueueIsEmpty
	}

	questionIndex := 0
	if game.Settings.ShuffleQuestions {
		randomNumber, _ := rand.Int(rand.Reader, big.NewInt(int64(len(unansweredQuestionIDs))))
		questionIndex = int(randomNumber.Int64())
	}

	return &unansweredQuestion{
		ID:    unansweredQuestionIDs[questionIndex],
		IsNew: true,
	}, nil
}
