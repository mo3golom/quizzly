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
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "orange", true),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "pink", true),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "amber", false),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "red", false),
		),
	), nil
}
