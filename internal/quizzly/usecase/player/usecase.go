package player

import (
	"context"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/player"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	maxNameLength = 25
)

type Usecase struct {
	players player.Repository
}

func NewUsecase(
	players player.Repository,
) contracts.PLayerUsecase {
	return &Usecase{
		players: players,
	}
}

func (u *Usecase) Create(ctx context.Context, in *model.Player) error {
	if in.ID == uuid.Nil {
		in.ID = uuid.New()
	}

	if utf8.RuneCountInString(in.Name) > 50 {
		in.Name = in.Name[:maxNameLength]
	}

	return u.players.Insert(ctx, in)
}

func (u *Usecase) Update(ctx context.Context, in *model.Player) error {
	return u.players.Update(ctx, in)
}

func (u *Usecase) Get(ctx context.Context, ids []uuid.UUID) ([]model.Player, error) {
	return u.players.GetByIDs(ctx, ids)
}

func (u *Usecase) GetByUsers(ctx context.Context, userIDs []uuid.UUID) ([]model.Player, error) {
	return u.players.GetByUserIDs(ctx, userIDs)
}
