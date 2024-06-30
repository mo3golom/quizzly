package game

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/handlers"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
)

type (
	PostStartData struct {
		GameID uuid.UUID `schema:"game-id"`
	}

	PostStartHandler struct {
		uc contracts.GameUsecase
	}
)

func NewPostStartHandler(uc contracts.GameUsecase) *PostStartHandler {
	return &PostStartHandler{uc: uc}
}

func (h *PostStartHandler) Handle(_ http.ResponseWriter, request *http.Request, in PostStartData) (templ.Component, error) {
	err := h.uc.Start(request.Context(), in.GameID)
	if err != nil {
		return nil, err
	}

	game, err := h.uc.Get(request.Context(), in.GameID)
	if err != nil {
		return nil, err
	}

	return frontendAdminGame.Header(
		&handlers.Game{
			ID:     game.ID,
			Status: game.Status,
			Link:   getGameLink(game.ID, request),
		},
	), nil
}
