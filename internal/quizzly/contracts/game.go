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
		Title    *string
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
