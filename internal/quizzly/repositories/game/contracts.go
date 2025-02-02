package game

import (
	"context"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/model"
)

type (
	Spec struct {
		IDs       []uuid.UUID
		AuthorID  *uuid.UUID
		IsPrivate *bool
		Limit     int64
		Statuses  []model.GameStatus
	}

	QuestionsSpec struct {
		IDs    []uuid.UUID
		GameID *uuid.UUID
	}

	Order struct {
		Field     string
		Direction string
	}

	Repository interface {
		Upsert(ctx context.Context, in *model.Game) error
		GetBySpec(ctx context.Context, spec *Spec) ([]model.Game, error)

		InsertQuestion(ctx context.Context, in *model.Question) error
		UpdateQuestion(ctx context.Context, in *model.Question) error
		DeleteQuestion(ctx context.Context, id uuid.UUID) error
		GetQuestionsBySpec(ctx context.Context, spec *QuestionsSpec) ([]model.Question, error)
	}
)
