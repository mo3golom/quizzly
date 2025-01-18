package question

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
)

type (
	GetListData struct {
		GameID uuid.UUID `schema:"game_id"`

		InContainer bool `schema:"in_container"`
	}

	GetHandler struct {
		service *service
	}
)

func NewGetHandler(uc contracts.GameUsecase) *GetHandler {
	return &GetHandler{
		service: &service{uc: uc},
	}
}

func (h *GetHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetListData) (templ.Component, error) {
	return h.service.list(request.Context(), in.GameID, in.InContainer)
}
