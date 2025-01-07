package contracts

import (
	"context"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/model"
)

type (
	AcceptAnswersIn struct {
		GameID     uuid.UUID
		PlayerID   uuid.UUID
		QuestionID uuid.UUID
		Answers    []string
	}

	AcceptAnswersOut struct {
		IsCorrect bool
		Details   []AnswerResult

		RightAnswers []model.AnswerOption
	}

	AnswerResult struct {
		Answer    string
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

	GetExtendedSessionsOut struct {
		Result     []model.ExtendedSession
		TotalCount int64
	}

	SessionUsecase interface {
		Start(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error
		Finish(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error
		Restart(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error

		AcceptAnswers(ctx context.Context, in *AcceptAnswersIn) (*AcceptAnswersOut, error)
		GetCurrentState(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*SessionState, error)

		GetStatistics(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*model.SessionStatistics, error)
		GetExtendedSessions(ctx context.Context, gameID uuid.UUID, page int64, limit int64) (*GetExtendedSessionsOut, error)
	}
)
