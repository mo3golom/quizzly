package contracts

import (
	"context"
	"quizzly/internal/quizzly/model"

	"github.com/google/uuid"
)

type (
	CreateGameIn struct {
		AuthorID uuid.UUID
		Type     model.GameType
		Settings model.GameSettings
	}

	GameUsecase interface {
		Create(ctx context.Context, in *CreateGameIn) (uuid.UUID, error)
		Start(ctx context.Context, id uuid.UUID) error
		Finish(ctx context.Context, id uuid.UUID) error
		Get(ctx context.Context, id uuid.UUID) (*model.Game, error)
		GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]model.Game, error)
		GetStatistics(ctx context.Context, id uuid.UUID) (*model.GameStatistics, error)

		AddQuestion(ctx context.Context, gameID uuid.UUID, questionID ...uuid.UUID) error
		GetQuestions(ctx context.Context, gameID uuid.UUID) ([]uuid.UUID, error)
	}

	AcceptAnswersIn struct {
		GameID     uuid.UUID
		PlayerID   uuid.UUID
		QuestionID uuid.UUID
		Answers    []model.AnswerOptionID
	}

	AcceptAnswersOut struct {
		IsCorrect bool
		Details   []AnswerResult
	}

	AnswerResult struct {
		Answer    model.AnswerOptionID
		IsCorrect bool
	}

	SessionState struct {
		Status          model.SessionStatus
		CurrentQuestion *model.Question
		Progress        Progress
	}

	Progress struct {
		Answered int64
		Total    int64
	}

	SessionUsecase interface {
		Start(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error
		Finish(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error

		AcceptAnswers(ctx context.Context, in *AcceptAnswersIn) (*AcceptAnswersOut, error)
		GetCurrentState(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*SessionState, error)
		GetStatistics(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*model.SessionStatistics, error)
		GetSessions(ctx context.Context, gameID uuid.UUID) ([]model.SessionExtended, error)
	}
)
