package session

import (
	"context"
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

type Usecase struct {
	sessions  session.Repository
	games     game.Repository
	questions question.Repository
	players   player.Repository
	template  transactional.Template
}

func NewUsecase(
	sessions session.Repository,
	games game.Repository,
	questions question.Repository,
	template transactional.Template,
) contracts.SessionUsecase {
	return &Usecase{
		sessions:  sessions,
		games:     games,
		questions: questions,
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

		if _, err := u.getActiveSession(ctx, tx, in.PlayerID, in.GameID); err != nil {
			return err
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

		specificQuestion, err := u.questions.Get(ctx, in.QuestionID)
		if err != nil {
			return err
		}

		result, err = checkAnswers(specificQuestion, in.Answers)
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

func (u *Usecase) NextQuestion(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*model.Question, error) {
	var result *model.Question
	return result, u.template.Execute(ctx, func(tx transactional.Tx) error {
		if _, err := u.getActiveGame(ctx, tx, gameID); err != nil {
			return err
		}

		specificSession, err := u.getActiveSession(ctx, tx, playerID, gameID)
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
		if slices.Contains(sessionItems, func(item model.SessionItem) bool {
			return item.AnsweredAt == nil
		}) {
			return contracts.ErrUnansweredQuestionExists
		}

		answeredQuestions := slices.Filter(sessionItems, func(item model.SessionItem) bool {
			return item.AnsweredAt != nil
		})
		unansweredQuestions, err := u.games.GetQuestionIDsBySpec(
			ctx,
			tx,
			&game.Spec{
				GameID: gameID,
				ExcludeQuestionIDs: slices.SafeMap(answeredQuestions, func(i model.SessionItem) uuid.UUID {
					return i.QuestionID
				}),
			},
		)

		if len(unansweredQuestions) == 0 {
			return contracts.ErrQuestionQueueIsEmpty
		}

		result, err = u.questions.Get(ctx, unansweredQuestions[0])
		if err != nil {
			return err
		}

		return u.sessions.InsertSessionItem(ctx, tx, &model.SessionItem{
			SessionID:  specificSession.ID,
			QuestionID: result.ID,
		})
	})
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

func (u *Usecase) getActiveSession(ctx context.Context, tx transactional.Tx, playerID uuid.UUID, gameID uuid.UUID) (*model.Session, error) {
	specificSession, err := u.sessions.GetBySpecWithTx(ctx, tx, &session.Spec{
		PlayerID: playerID,
		GameID:   gameID,
	})
	if err != nil {
		return nil, err
	}
	if specificSession.Status != model.SessionStatusStarted {
		return nil, errors.New("player session isn't active")
	}

	return specificSession, nil
}
