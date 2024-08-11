package question

import (
	"context"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type (
	Spec struct {
		IDs      []uuid.UUID
		AuthorID *uuid.UUID
		Page     *Page
	}

	Page struct {
		Number int64
		Limit  int64
	}

	GetBySpecOut struct {
		Result     []model.Question
		TotalCount int64
	}

	Repository interface {
		Insert(ctx context.Context, tx transactional.Tx, in *model.Question) error
		Delete(ctx context.Context, tx transactional.Tx, id uuid.UUID) error
		GetBySpec(ctx context.Context, spec *Spec) (*GetBySpecOut, error)
	}
)
