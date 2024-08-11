package contracts

import (
	"context"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/model"
)

type (
	GetByAuthorOut struct {
		Result     []model.Question
		TotalCount int64
	}

	QuestionUsecase interface {
		Create(ctx context.Context, in *model.Question) error
		Delete(ctx context.Context, id uuid.UUID) error
		GetByAuthor(ctx context.Context, authorID uuid.UUID, page int64, limit int64) (*GetByAuthorOut, error)
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Question, error)
	}
)
