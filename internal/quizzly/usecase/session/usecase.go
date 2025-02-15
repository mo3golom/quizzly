package session

import (
	"context"
	"errors"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/player"
	"quizzly/internal/quizzly/repositories/session"
)

type (
	AnswerOptionIDAcceptor interface {
		Accept(question *model.Question, answers []string) (*contracts.AcceptAnswersOut, error)
	}

	Usecase struct {
		sessions session.Repository
		games    game.Repository
		players  player.Repository
		trm      trm.Manager

		optionIDAcceptors map[model.QuestionType]AnswerOptionIDAcceptor
	}
)

func NewUsecase(
	sessions session.Repository,
	games game.Repository,
	players player.Repository,
	trm trm.Manager,
	optionIDAcceptors map[model.QuestionType]AnswerOptionIDAcceptor,
) contracts.SessionUsecase {
	return &Usecase{
		sessions:          sessions,
		games:             games,
		players:           players,
		trm:               trm,
		optionIDAcceptors: optionIDAcceptors,
	}
}

func (u *Usecase) Start(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error {
	return u.trm.Do(ctx, func(ctx context.Context) error {
		if _, err := u.getActiveGame(ctx, gameID); err != nil {
			return err
		}

		specificPlayers, err := u.players.GetByIDs(ctx, []uuid.UUID{playerID})
		if err != nil {
			return err
		}
		if len(specificPlayers) == 0 {
			err = u.players.Insert(ctx, &model.Player{
				ID:   playerID,
				Name: "unknown",
			})
			if err != nil {
				return err
			}
		}

		return u.sessions.Insert(
			ctx,
			&model.Session{
				PlayerID: playerID,
				GameID:   gameID,
				Status:   model.SessionStatusStarted,
			},
		)
	})
}

func (u *Usecase) Finish(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error {
	return u.trm.Do(ctx, func(ctx context.Context) error {
		if _, err := u.getActiveGame(ctx, gameID); err != nil {
			return err
		}

		specificPlayerGame, err := u.sessions.GetBySpec(ctx, &session.Spec{
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
		return u.sessions.Update(ctx, specificPlayerGame)
	})
}

func (u *Usecase) Restart(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error {
	return u.trm.Do(ctx, func(ctx context.Context) error {
		if _, err := u.getActiveGame(ctx, gameID); err != nil {
			return err
		}

		specificPlayerGame, err := u.sessions.GetBySpec(ctx, &session.Spec{
			PlayerID: playerID,
			GameID:   gameID,
		})
		if errors.Is(err, contracts.ErrSessionNotFound) {
			return nil
		}

		if err != nil {
			return err
		}

		err = u.sessions.DeleteSessionItemsBySessionID(ctx, specificPlayerGame.ID)
		if err != nil {
			return err
		}

		specificPlayerGame.Status = model.SessionStatusStarted
		return u.sessions.Update(ctx, specificPlayerGame)
	})
}

func (u *Usecase) GetStatistics(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*model.SessionStatistics, error) {
	var result *model.SessionStatistics
	return result, u.trm.Do(ctx, func(ctx context.Context) error {
		specificPlayerGame, err := u.sessions.GetBySpec(ctx, &session.Spec{
			PlayerID: playerID,
			GameID:   gameID,
		})
		if err != nil {
			return err
		}
		if specificPlayerGame.Status != model.SessionStatusFinished {
			return contracts.ErrSessionNotFinished
		}

		sessionItems, err := u.sessions.GetSessionBySpec(
			ctx,
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

func (u *Usecase) getActiveGame(ctx context.Context, gameID uuid.UUID) (*model.Game, error) {
	specificGames, err := u.games.GetBySpec(ctx, &game.Spec{
		IDs: []uuid.UUID{gameID},
	})
	if err != nil {
		return nil, err
	}
	if len(specificGames) == 0 {
		return nil, contracts.ErrGameNotFound
	}

	specificGame := specificGames[0]
	if specificGame.Status != model.GameStatusStarted {
		return nil, errors.New("game isn't started")
	}

	return &specificGame, nil
}

func (u *Usecase) getSession(ctx context.Context, playerID uuid.UUID, gameID uuid.UUID) (*model.Session, error) {
	specificSession, err := u.sessions.GetBySpec(ctx, &session.Spec{
		PlayerID: playerID,
		GameID:   gameID,
	})
	if errors.Is(err, contracts.ErrSessionNotFound) {
		err := u.Start(ctx, gameID, playerID)
		if err != nil {
			return nil, err
		}

		return u.sessions.GetBySpec(ctx, &session.Spec{
			PlayerID: playerID,
			GameID:   gameID,
		})
	}
	if err != nil {
		return nil, err
	}

	return specificSession, nil
}
