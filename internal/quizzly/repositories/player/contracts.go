package player

import (
	"context"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type (
	Repository interface {
		Insert(ctx context.Context, tx transactional.Tx, in *model.Player) error
		Get(ctx context.Context, id uuid.UUID) (*model.Player, error)
	}
)
