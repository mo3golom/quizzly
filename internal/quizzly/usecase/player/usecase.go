package player

import (
	"context"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/player"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type Usecase struct {
	players  player.Repository
	template transactional.Template
}

func NewUsecase(
	players player.Repository,
	template transactional.Template,
) contracts.PLayerUsecase {
	return &Usecase{
		players:  players,
		template: template,
	}
}

func (u *Usecase) Create(ctx context.Context, in *model.Player) error {
	if in.ID == uuid.Nil {
		in.ID = uuid.New()
	}

	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.players.Insert(ctx, tx, in)
	})
}

func (u *Usecase) Get(ctx context.Context, ids []uuid.UUID) ([]model.Player, error) {
	return u.players.GetByIDs(ctx, ids)
}
