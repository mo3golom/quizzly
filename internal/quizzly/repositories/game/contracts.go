package game

import (
	"context"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
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
		Order  *Order
	}

	Order struct {
		Field     string
		Direction string
	}

	Repository interface {
		Upsert(ctx context.Context, tx transactional.Tx, in *model.Game) error
		GetBySpec(ctx context.Context, spec *Spec) ([]model.Game, error)
		GetBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *Spec) ([]model.Game, error)

		InsertQuestion(ctx context.Context, tx transactional.Tx, in *model.Question) error
		UpdateQuestion(ctx context.Context, tx transactional.Tx, in *model.Question) error
		DeleteQuestion(ctx context.Context, tx transactional.Tx, id uuid.UUID) error
		GetQuestionsBySpec(ctx context.Context, spec *QuestionsSpec) ([]model.Question, error)
		GetQuestionsBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *QuestionsSpec) ([]model.Question, error)
	}
)
