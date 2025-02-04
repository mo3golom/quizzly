package game

import (
	"fmt"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/handlers"
)

func convertModelGameToHandlersGame(game *model.Game) *handlers.Game {
	title := fmt.Sprintf("Игра без названия от %s", game.CreatedAt.Format("02.01.2006"))
	if game.Title != nil {
		title = *game.Title
	}

	return &handlers.Game{
		ID:        game.ID,
		Status:    game.Status,
		Title:     title,
		CreatedAt: game.CreatedAt,
		Settings: handlers.GameSettings{
			ShuffleQuestions: game.Settings.ShuffleQuestions,
			ShuffleAnswers:   game.Settings.ShuffleAnswers,
			ShowRightAnswers: game.Settings.ShowRightAnswers,
			InputCustomName:  game.Settings.InputCustomName,
			IsPrivate:        game.Settings.IsPrivate,
		},
	}
}
