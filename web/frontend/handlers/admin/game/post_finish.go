package game

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/services/link"
)

type (
	PostFinishData struct {
		GameID uuid.UUID `schema:"game-id"`
	}

	PostFinishHandler struct {
		uc contracts.GameUsecase

		service *service
	}
)

func NewPostFinishHandler(
	uc contracts.GameUsecase,
	linkService link.Service,
) *PostFinishHandler {
	return &PostFinishHandler{
		uc: uc,
		service: &service{
			uc:          uc,
			linkService: linkService,
		},
	}
}

func (h *PostFinishHandler) Handle(_ http.ResponseWriter, request *http.Request, in PostFinishData) (templ.Component, error) {
	err := h.uc.Finish(request.Context(), in.GameID)
	if err != nil {
		return nil, err
	}

	return h.service.getGamePage(request, in.GameID)
}
