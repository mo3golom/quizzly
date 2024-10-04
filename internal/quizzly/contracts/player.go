package contracts

import (
	"context"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/model"
)

type (
	PLayerUsecase interface {
		Create(ctx context.Context, in *model.Player) error
		Update(ctx context.Context, in *model.Player) error
		Get(ctx context.Context, ids []uuid.UUID) ([]model.Player, error)
		GetByUsers(ctx context.Context, userIDs []uuid.UUID) ([]model.Player, error)
	}
)
