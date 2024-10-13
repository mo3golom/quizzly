package game

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/services/page"
	"quizzly/web/frontend/services/player"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

const (
	getRestartPageTitle = "Перезапуск игры"
)

type (
	GetRestartPageData struct {
		GameID *uuid.UUID `schema:"id"`
	}

	GetRestartPageHandler struct {
		gameUC    contracts.GameUsecase
		sessionUC contracts.SessionUsecase

		playerService player.Service
	}
)

func NewGetRestartPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerService player.Service,
) *GetRestartPageHandler {
	return &GetRestartPageHandler{
		gameUC:        gameUC,
		sessionUC:     sessionUC,
		playerService: playerService,
	}
}

func (h *GetRestartPageHandler) Handle(writer http.ResponseWriter, request *http.Request, in GetRestartPageData) (templ.Component, error) {
	gameID := in.GameID
	if pathGameID := request.PathValue(pathValueGameID); pathGameID != "" {
		tempGameID, err := uuid.Parse(pathGameID)
		if err != nil {
			return nil, err
		}

		gameID = &tempGameID
	}

	game, err := h.gameUC.Get(request.Context(), *gameID)
	if errors.Is(err, contracts.ErrGameNotFound) {
		return page.PublicIndexPage(
			request.Context(),
			getRestartPageTitle,
			frontendPublicGame.StartPage("Игра не найдена."),
		), nil
	}
	if err != nil {
		return nil, err
	}
	if game.Status == model.GameStatusFinished {
		return page.PublicIndexPage(
			request.Context(),
			getRestartPageTitle,
			frontendPublicGame.StartPage("Игра уже завершена."),
		), nil
	}
	if game.Status == model.GameStatusCreated {
		return page.PublicIndexPage(
			request.Context(),
			getRestartPageTitle,
			frontendPublicGame.StartPage("Игра еще не началась. Подождите немного или попросите автора запустить игру."),
		), nil
	}

	currentPlayer, err := h.playerService.GetPlayer(writer, request, game.ID)
	if err != nil {
		return nil, err
	}

	err = h.sessionUC.Restart(request.Context(), game.ID, currentPlayer.ID)
	if err != nil {
		return nil, err
	}

	return frontendComponents.Redirect(gameLink(game.ID)), nil
}
