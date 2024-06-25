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
	}

	Repository interface {
		Insert(ctx context.Context, tx transactional.Tx, in *model.Question) error
		GetBySpec(ctx context.Context, spec *Spec) ([]model.Question, error)
	}
)
