package question

import (
	"net/http"
	"quizzly/internal/quizzly/contracts"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type (
	GetDeleteData struct {
		ID     uuid.UUID `schema:"id"`
		GameID uuid.UUID `schema:"game_id"`
	}

	GetDeleteHandler struct {
		uc      contracts.GameUsecase
		service *service
	}
)

func NewPostDeleteHandler(uc contracts.GameUsecase) *GetDeleteHandler {
	return &GetDeleteHandler{
		uc:      uc,
		service: &service{uc: uc},
	}
}

func (h *GetDeleteHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetDeleteData) (templ.Component, error) {
	err := h.uc.DeleteQuestion(request.Context(), in.ID)
	if err != nil {
		return nil, err
	}

	return h.service.list(request.Context(), in.GameID, true, true)
}
