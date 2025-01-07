package game

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/services/link"
)

type (
	PostStartData struct {
		GameID uuid.UUID `schema:"game-id"`
	}

	PostStartHandler struct {
		uc contracts.GameUsecase

		service *service
	}
)

func NewPostStartHandler(
	uc contracts.GameUsecase,
	linkService link.Service,
) *PostStartHandler {
	return &PostStartHandler{
		uc: uc,
		service: &service{
			uc:          uc,
			linkService: linkService,
		},
	}
}

func (h *PostStartHandler) Handle(_ http.ResponseWriter, request *http.Request, in PostStartData) (templ.Component, error) {
	err := h.uc.Start(request.Context(), in.GameID)
	if err != nil {
		return nil, err
	}

	return h.service.getGamePage(request, in.GameID)
}
