package player

import (
	"context"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/model"
)

type (
	Repository interface {
		Insert(ctx context.Context, in *model.Player) error
		Update(ctx context.Context, in *model.Player) error
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Player, error)
		GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]model.Player, error)
	}
)
