package admin

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/handlers"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
)

type (
	PostFinishData struct {
		GameID uuid.UUID `schema:"game-id"`
	}

	PostFinishHandler struct {
		uc contracts.GameUsecase
	}
)

func NewPostFinishHandler(uc contracts.GameUsecase) *PostFinishHandler {
	return &PostFinishHandler{uc: uc}
}

func (h *PostFinishHandler) Handle(_ http.ResponseWriter, request *http.Request, in PostFinishData) (templ.Component, error) {
	err := h.uc.Finish(request.Context(), in.GameID)
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
