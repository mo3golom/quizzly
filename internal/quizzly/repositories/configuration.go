package repositories

import (
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/player"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs"

	"github.com/jmoiron/sqlx"
)

type (
	Configuration struct {
		Game     structs.Singleton[game.Repository]
		Session  structs.Singleton[session.Repository]
		Question structs.Singleton[question.Repository]
		Player   structs.Singleton[player.Repository]
	}
)

func NewConfiguration(db *sqlx.DB) *Configuration {
	return &Configuration{
		Game: structs.NewSingleton(func() (game.Repository, error) {
			return game.NewRepository(), nil
		}),
		Session: structs.NewSingleton(func() (session.Repository, error) {
			return session.NewRepository(db), nil
		}),
		Question: structs.NewSingleton(func() (question.Repository, error) {
			return question.NewRepository(db), nil
		}),
		Player: structs.NewSingleton(func() (player.Repository, error) {
			return player.NewRepository(db), nil
		}),
	}
}
