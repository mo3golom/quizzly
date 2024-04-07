package contracts

import (
	"context"
	"quizzly/internal/quizzly/model"

	"github.com/google/uuid"
)

type (
	PlayerUsecase interface {
		Create(ctx context.Context, data model.Player) error
		Get(ctx context.Context, id uuid.UUID) (*model.Player, error)
	}
)
