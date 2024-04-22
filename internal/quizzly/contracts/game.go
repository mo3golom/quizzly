package contracts

import (
	"context"
	"quizzly/internal/quizzly/model"

	"github.com/google/uuid"
)

type (
	CreateGameIn struct {
		Type     model.GameType
		Settings model.GameSettings
	}

	GameUsecase interface {
		Create(ctx context.Context, in *CreateGameIn) (uuid.UUID, error)
		Start(ctx context.Context, id uuid.UUID) error
		Finish(ctx context.Context, id uuid.UUID) error

		AddQuestion(ctx context.Context, gameID uuid.UUID, questionID uuid.UUID) error
	}

	AcceptAnswersIn struct {
		GameID     uuid.UUID
		PlayerID   uuid.UUID
		QuestionID uuid.UUID
		Answers    []string
	}

	AcceptAnswersOut struct {
		IsCorrect bool
		Details   []AnswerResult
	}

	AnswerResult struct {
		Answer    string
		IsCorrect bool
	}

	SessionUsecase interface {
		Start(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error
		Finish(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error

		AcceptAnswers(ctx context.Context, in *AcceptAnswersIn) (*AcceptAnswersOut, error)
		NextQuestion(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*model.Question, error)
	}
)
