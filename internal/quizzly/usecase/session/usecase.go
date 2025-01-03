package session

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/player"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/transactional"
)

type (
	unansweredQuestion struct {
		ID    uuid.UUID
		IsNew bool
	}

	AnswerOptionIDAcceptor interface {
		Accept(question *model.Question, answers []model.AnswerOptionID) (*contracts.AcceptAnswersOut, error)
	}

	AnswerTextAcceptor interface {
		Accept(question *model.Question, answers []string) (*contracts.AcceptAnswersOut, error)
	}

	Usecase struct {
		sessions  session.Repository
		games     game.Repository
		questions question.Repository
		players   player.Repository
		template  transactional.Template

		optionIDAcceptors map[model.QuestionType]AnswerOptionIDAcceptor
		textAcceptors     map[model.QuestionType]AnswerTextAcceptor
	}
)

func NewUsecase(
	sessions session.Repository,
	games game.Repository,
	questions question.Repository,
	players player.Repository,
	template transactional.Template,
	optionIDAcceptors map[model.QuestionType]AnswerOptionIDAcceptor,
	textAcceptors map[model.QuestionType]AnswerTextAcceptor,
) contracts.SessionUsecase {
	return &Usecase{
		sessions:          sessions,
		games:             games,
		questions:         questions,
		players:           players,
		template:          template,
		optionIDAcceptors: optionIDAcceptors,
		textAcceptors:     textAcceptors,
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
		if errors.Is(err, contracts.ErrSessionNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		specificPlayerGame.Status = model.SessionStatusFinished
		return u.sessions.Update(ctx, tx, specificPlayerGame)
	})
}

func (u *Usecase) Restart(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		if _, err := u.getActiveGame(ctx, tx, gameID); err != nil {
			return err
		}

		specificPlayerGame, err := u.sessions.GetBySpecWithTx(ctx, tx, &session.Spec{
			PlayerID: playerID,
			GameID:   gameID,
		})
		if errors.Is(err, contracts.ErrSessionNotFound) {
			return nil
		}

		if err != nil {
			return err
		}

		err = u.sessions.DeleteSessionItemsBySessionID(ctx, tx, specificPlayerGame.ID)
		if err != nil {
			return err
		}

		specificPlayerGame.Status = model.SessionStatusStarted
		return u.sessions.Update(ctx, tx, specificPlayerGame)
	})
}

func (u *Usecase) GetStatistics(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*model.SessionStatistics, error) {
	var result *model.SessionStatistics
	return result, u.template.Execute(ctx, func(tx transactional.Tx) error {
		specificPlayerGame, err := u.sessions.GetBySpecWithTx(ctx, tx, &session.Spec{
			PlayerID: playerID,
			GameID:   gameID,
		})
		if err != nil {
			return err
		}
		if specificPlayerGame.Status != model.SessionStatusFinished {
			return contracts.ErrSessionNotFinished
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

func (u *Usecase) GetExtendedSessions(ctx context.Context, gameID uuid.UUID, page int64, limit int64) (*contracts.GetExtendedSessionsOut, error) {
	sessions, err := u.sessions.GetExtendedSessionsBySpec(ctx, &session.GetExtendedSessionSpec{
		GameID: gameID,
		Page: &session.Page{
			Number: page,
			Limit:  limit,
		},
	})
	if err != nil {
		return nil, err
	}

	return &contracts.GetExtendedSessionsOut{
		Result:     sessions.Result,
		TotalCount: sessions.TotalCount,
	}, nil
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
	if errors.Is(err, contracts.ErrSessionNotFound) {
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
