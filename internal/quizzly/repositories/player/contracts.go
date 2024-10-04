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
		Update(ctx context.Context, tx transactional.Tx, in *model.Player) error
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Player, error)
		GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]model.Player, error)
	}
)
