package question

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/pkg/structs"
	"quizzly/web/frontend/services/question"
)

type (
	GetHandler struct {
		service question.Service
	}
)

func NewGetHandler(service question.Service) *GetHandler {
	return &GetHandler{service: service}
}

func (h *GetHandler) Handle(_ http.ResponseWriter, request *http.Request, _ struct{}) (templ.Component, error) {
	return h.service.List(
		request.Context(),
		&question.Spec{
			AuthorID: structs.Pointer(uuid.New()),
		},
		&question.ListOptions{},
	)
}
