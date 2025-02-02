package question

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/model"
	frontendAdminQuestion "quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
)

type (
	GetFormData struct {
		GameID uuid.UUID `schema:"game_id"`
	}

	GetFormHandler struct {
	}
)

func NewGetFormHandler() *GetFormHandler {
	return &GetFormHandler{}
}

func (h *GetFormHandler) Handle(_ http.ResponseWriter, _ *http.Request, in GetFormData) (templ.Component, error) {
	return frontendAdminQuestion.Form(
		in.GameID,
		model.QuestionTypeChoice,
		frontendComponents.Composition(
			frontendAdminQuestion.QuestionImageInput(),
			frontendAdminQuestion.QuestionTextInput(),
		),
		frontendComponents.Composition(
			frontendAdminQuestion.AnswerChoiceInput(0, uuid.New(), true),
			frontendAdminQuestion.AnswerChoiceInput(1, uuid.New(), true),
			frontendAdminQuestion.AnswerChoiceInput(2, uuid.New(), false),
			frontendAdminQuestion.AnswerChoiceInput(3, uuid.New(), false),
		),
	), nil
}
