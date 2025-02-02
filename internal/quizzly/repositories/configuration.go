package repositories

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/player"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs"

	"github.com/jmoiron/sqlx"
)

type (
	Configuration struct {
		Game    structs.Singleton[game.Repository]
		Session structs.Singleton[session.Repository]
		Player  structs.Singleton[player.Repository]
	}
)

func NewConfiguration(db *sqlx.DB) *Configuration {
	trmsqlxGetter := trmsqlx.DefaultCtxGetter

	return &Configuration{
		Game: structs.NewSingleton(func() (game.Repository, error) {
			return game.NewRepository(db, trmsqlxGetter), nil
		}),
		Session: structs.NewSingleton(func() (session.Repository, error) {
			return session.NewRepository(db, trmsqlxGetter), nil
		}),
		Player: structs.NewSingleton(func() (player.Repository, error) {
			return player.NewRepository(db, trmsqlxGetter), nil
		}),
	}
}
