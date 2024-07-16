package question

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	frontendAdminQuestion "quizzly/web/frontend/templ/admin/question"
)

type (
	GetDeleteData struct {
		ID uuid.UUID `schema:"id"`
	}

	GetDeleteHandler struct {
		uc contracts.QuestionUsecase
	}
)

func NewPostDeleteHandler(uc contracts.QuestionUsecase) *GetDeleteHandler {
	return &GetDeleteHandler{
		uc: uc,
	}
}

func (h *GetDeleteHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetDeleteData) (templ.Component, error) {
	questions, err := h.uc.GetByIDs(request.Context(), []uuid.UUID{in.ID})
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, errors.New("question not found")
	}

	err = h.uc.Delete(request.Context(), in.ID)
	if err != nil {
		return nil, err
	}

	question := questions[0]
	return frontendAdminQuestion.QuestionListItem(
		question.ID,
		question.Text,
		question.Type,
		nil,
		frontendAdminQuestion.Options{
			WithActions:         true,
			WithDisabledOverlay: true,
		},
	), nil
}
