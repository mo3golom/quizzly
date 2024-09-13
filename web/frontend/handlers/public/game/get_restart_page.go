package game

import (
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	frontend "quizzly/web/frontend/templ"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

const (
	getRestartPageTitle = "Перезапуск игры"
)

type (
	GetRestartPageData struct {
		GameID uuid.UUID `schema:"id"`
	}

	GetRestartPageHandler struct {
		gameUC    contracts.GameUsecase
		sessionUC contracts.SessionUsecase
	}
)

func NewGetRestartPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
) *GetRestartPageHandler {
	return &GetRestartPageHandler{
		gameUC:    gameUC,
		sessionUC: sessionUC,
	}
}

func (h *GetRestartPageHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetRestartPageData) (templ.Component, error) {
	game, err := h.gameUC.Get(request.Context(), in.GameID)
	if errors.Is(err, contracts.ErrGameNotFound) {
		return frontend.PublicPageComponent(
			getRestartPageTitle,
			frontendPublicGame.StartPage("Игра не найдена."),
		), nil
	}
	if err != nil {
		return nil, err
	}
	if game.Status == model.GameStatusFinished {
		return frontend.PublicPageComponent(
			getRestartPageTitle,
			frontendPublicGame.StartPage("Игра уже завершена."),
		), nil
	}
	if game.Status == model.GameStatusCreated {
		return frontend.PublicPageComponent(
			getRestartPageTitle,
			frontendPublicGame.StartPage("Игра еще не началась. Подождите немного или попросите автора запустить игру."),
		), nil
	}

	playerID := getPlayerID(request)
	if playerID == uuid.Nil {
		return frontendComponents.Redirect(h.getGameLink(game.ID)), nil
	}

	err = h.sessionUC.Restart(request.Context(), game.ID, playerID)
	if err != nil {
		return nil, err
	}

	return frontendComponents.Redirect(h.getGameLink(game.ID)), nil
}

func (h *GetRestartPageHandler) getGameLink(gameID uuid.UUID) string {
	return fmt.Sprintf("/game/play?id=%s", gameID.String())
}
