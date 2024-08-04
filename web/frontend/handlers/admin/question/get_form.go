package question

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/model"
	frontend "quizzly/web/frontend/templ"
	frontendAdminQuestion "quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
)

const (
	formAddTitle = "Добавление вопроса"
	listUrl      = "/question/list"
)

type (
	GetFormData struct{}

	GetFormHandler struct {
	}
)

func NewGetFormHandler() *GetFormHandler {
	return &GetFormHandler{}
}

func (h *GetFormHandler) Handle(_ http.ResponseWriter, _ *http.Request, _ GetFormData) (templ.Component, error) {
	return frontend.AdminPageComponent(
		formAddTitle,
		frontendComponents.Composition(
			frontendComponents.BackLink(listUrl),
			frontendComponents.Tabs(
				uuid.New(),
				frontendComponents.Tab{
					Name:    "Один ответ",
					Content: singleChoiceForm(),
				},
				frontendComponents.Tab{
					Name:    "Несколько ответов",
					Content: multipleChoiceForm(),
				},
				frontendComponents.Tab{
					Name:    "Ввод слова",
					Content: fillTheGapForm(),
				},
			),
		),
	), nil
}

func singleChoiceForm() templ.Component {
	return frontendAdminQuestion.Form(
		model.QuestionTypeChoice,
		frontendComponents.Composition(
			frontendAdminQuestion.QuestionImageInput(),
			frontendAdminQuestion.QuestionTextInput(),
		),
		frontendComponents.Composition(
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "orange", false, true),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "pink", false, true),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "amber", false, false),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "red", false, false),
		),
	)
}

func multipleChoiceForm() templ.Component {
	return frontendAdminQuestion.Form(
		model.QuestionTypeMultipleChoice,
		frontendComponents.Composition(
			frontendAdminQuestion.QuestionImageInput(),
			frontendAdminQuestion.QuestionTextInput(),
			frontendAdminQuestion.QuestionMultipleChoiceOption(),
		),
		frontendComponents.Composition(
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "orange", true, true),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "pink", true, true),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "amber", true, false),
			frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "red", true, false),
		),
	)
}

func fillTheGapForm() templ.Component {
	return frontendAdminQuestion.Form(
		model.QuestionTypeFillTheGap,
		frontendComponents.Composition(
			frontendAdminQuestion.QuestionImageInput(),
			frontendAdminQuestion.QuestionTextInput(),
		),
		frontendComponents.Composition(
			frontendAdminQuestion.AnswerTextInput(),
		),
	)
}
