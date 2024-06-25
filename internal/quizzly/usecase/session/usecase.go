package session

import (
	"context"
	"database/sql"
	"errors"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/player"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"
	"time"

	"github.com/google/uuid"
)

type (
	unansweredQuestion struct {
		ID    uuid.UUID
		IsNew bool
	}

	Usecase struct {
		sessions  session.Repository
		games     game.Repository
		questions question.Repository
		players   player.Repository
		template  transactional.Template
	}
)

func NewUsecase(
	sessions session.Repository,
	games game.Repository,
	questions question.Repository,
	players player.Repository,
	template transactional.Template,
) contracts.SessionUsecase {
	return &Usecase{
		sessions:  sessions,
		games:     games,
		questions: questions,
		players:   players,
		template:  template,
	}
}

func (u *Usecase) Start(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		if _, err := u.getActiveGame(ctx, tx, gameID); err != nil {
			return err
		}

		specificPlayer, err := u.players.Get(ctx, playerID)
		if err != nil {
			return err
		}
		if specificPlayer == nil {
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

func (u *Usecase) AcceptAnswers(ctx context.Context, in *contracts.AcceptAnswersIn) (*contracts.AcceptAnswersOut, error) {
	var result *contracts.AcceptAnswersOut
	return result, u.template.Execute(ctx, func(tx transactional.Tx) error {
		if _, err := u.getActiveGame(ctx, tx, in.GameID); err != nil {
			return err
		}

		specificSession, err := u.getSession(ctx, tx, in.PlayerID, in.GameID)
		if err != nil {
			return err
		}
		if specificSession.Status != model.SessionStatusStarted {
			return contracts.ErrNotActiveSessionNotFound
		}

		specificPlayerSessionItems, err := u.sessions.GetSessionBySpecWithTx(ctx, tx, &session.ItemSpec{
			PlayerID:   in.PlayerID,
			GameID:     in.GameID,
			QuestionID: &in.QuestionID,
		})
		if err != nil {
			return err
		}
		if len(specificPlayerSessionItems) == 0 {
			return errors.New("player session is empty")
		}

		specificPlayerSessionItem := specificPlayerSessionItems[0]
		if specificPlayerSessionItem.AnsweredAt != nil {
			return errors.New("question is already answered")
		}

		specificQuestions, err := u.questions.GetBySpec(ctx, &question.Spec{
			IDs: []uuid.UUID{in.QuestionID},
		})
		if err != nil {
			return err
		}
		if len(specificQuestions) == 0 {
			return errors.New("question not found")
		}

		result, err = checkAnswers(&specificQuestions[0], in.Answers)
		if err != nil {
			return err
		}

		specificPlayerSessionItem.IsCorrect = structs.Pointer(result.IsCorrect)
		specificPlayerSessionItem.Answers = in.Answers
		specificPlayerSessionItem.AnsweredAt = structs.Pointer(time.Now())
		return u.sessions.UpdateSessionItem(
			ctx,
			tx,
			&specificPlayerSessionItem,
		)
	})
}

func (u *Usecase) GetCurrentState(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*contracts.SessionState, error) {
	var result *contracts.SessionState
	return result, u.template.Execute(ctx, func(tx transactional.Tx) error {
		if _, err := u.getActiveGame(ctx, tx, gameID); err != nil {
			return err
		}

		specificSession, err := u.getSession(ctx, tx, playerID, gameID)
		if err != nil {
			return err
		}
		if specificSession.Status == model.SessionStatusFinished {
			result = &contracts.SessionState{
				Status: specificSession.Status,
			}
			return nil
		}

		gameQuestions, err := u.games.GetQuestionIDsBySpec(
			ctx,
			tx,
			&game.Spec{
				GameID: gameID,
			},
		)
		if err != nil {
			return err
		}

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

		specificUnansweredQuestion, err := getUnansweredQuestionID(gameQuestions, sessionItems)
		if err != nil {
			return err
		}

		specificQuestions, err := u.questions.GetBySpec(ctx, &question.Spec{
			IDs: []uuid.UUID{specificUnansweredQuestion.ID},
		})
		if err != nil {
			return err
		}
		if len(specificQuestions) == 0 {
			return errors.New("question not found")
		}

		result = &contracts.SessionState{
			CurrentQuestion: &specificQuestions[0],
			Progress: contracts.Progress{
				Total: int64(len(gameQuestions)),
				Answered: int64(len(
					slices.Filter(sessionItems, func(item model.SessionItem) bool {
						return item.AnsweredAt != nil
					}),
				)),
			},
		}

		if !specificUnansweredQuestion.IsNew {
			return nil
		}
		return u.sessions.InsertSessionItem(ctx, tx, &model.SessionItem{
			SessionID:  specificSession.ID,
			QuestionID: result.CurrentQuestion.ID,
		})
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

func getUnansweredQuestionID(gameQuestions []uuid.UUID, sessionItems []model.SessionItem) (*unansweredQuestion, error) {
	unAnsweredSessionItems := slices.Filter(sessionItems, func(item model.SessionItem) bool {
		return item.AnsweredAt == nil
	})
	if len(unAnsweredSessionItems) > 0 {
		return &unansweredQuestion{
			ID:    unAnsweredSessionItems[0].QuestionID,
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

	return &unansweredQuestion{
		ID:    unansweredQuestionIDs[0],
		IsNew: true,
	}, nil
}
