package contracts

import (
	"context"
	"quizzly/internal/quizzly/model"
)

type (
	PLayerUsecase interface {
		Create(ctx context.Context, in *model.Player) error
	}
)
