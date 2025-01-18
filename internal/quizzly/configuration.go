package quizzly

import (
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories"
	"quizzly/internal/quizzly/usecase/game"
	"quizzly/internal/quizzly/usecase/player"
	"quizzly/internal/quizzly/usecase/session"
	"quizzly/internal/quizzly/usecase/session/acceptor"
	"quizzly/pkg/structs"
	"quizzly/pkg/transactional"

	"github.com/jmoiron/sqlx"
)

type (
	Configuration struct {
		Game    structs.Singleton[contracts.GameUsecase]
		Session structs.Singleton[contracts.SessionUsecase]
		Player  structs.Singleton[contracts.PLayerUsecase]
	}
)

func NewConfiguration(
	db *sqlx.DB,
	template transactional.Template,
) *Configuration {
	repos := repositories.NewConfiguration(db)

	return &Configuration{
		Game: structs.NewSingleton(func() (contracts.GameUsecase, error) {
			return game.NewUsecase(
				repos.Game.MustGet(),
				repos.Session.MustGet(),
				template,
			), nil
		}),
		Session: structs.NewSingleton(func() (contracts.SessionUsecase, error) {
			return session.NewUsecase(
				repos.Session.MustGet(),
				repos.Game.MustGet(),
				repos.Player.MustGet(),
				template,
				map[model.QuestionType]session.AnswerOptionIDAcceptor{
					model.QuestionTypeChoice:         acceptor.NewSingleChoiceAcceptor(),
					model.QuestionTypeOneOfChoice:    acceptor.NewOneOfChoiceAcceptor(),
					model.QuestionTypeMultipleChoice: acceptor.NewMultipleChoiceAcceptor(),
				},
			), nil
		}),
		Player: structs.NewSingleton(func() (contracts.PLayerUsecase, error) {
			return player.NewUsecase(
				repos.Player.MustGet(),
				template,
			), nil
		}),
	}
}
