package contracts

import (
	"context"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/model"
)

type (
	QuestionUsecase interface {
		Create(ctx context.Context, in *model.Question) error
		Delete(ctx context.Context, id uuid.UUID) error
		GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]model.Question, error)
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Question, error)
	}
)
