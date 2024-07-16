package question

import (
	"cmp"
	"fmt"
	"github.com/a-h/templ"
	"net/http"
	"quizzly/internal/quizzly/model"
	frontend "quizzly/web/frontend/templ"
	"quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
	goSlices "slices"
)

const (
	formAddTitle = "Добавление вопроса"
	submitUrl    = "/question"
	listUrl      = "/game/list"
)

const (
	paramAnswerOptionType = "answer_option_type"
	paramBGColor          = "bg_color"
	paramQuestionType     = "question_type"
)

const (
	answerOptionFormTypeForm = "a_o_f"
	answerOptionFormTypeAdd  = "a_o_a"
)

var (
	questionTypes = map[model.QuestionType]questionTypeItem{
		model.QuestionTypeChoice: {
			name:  "Один ответ",
			color: "blue-500",
		},
		model.QuestionTypeMultipleChoice: {
			name:  "Несколько ответов",
			color: "amber-500",
		},
	}
)

type (
	questionTypeItem struct {
		name  string
		color string
	}

	GetFormData struct {
		ActiveQuestionType   string `schema:"question_type"`
		AnswerOptionFormType string `schema:"answer_option_type"`
		BGColor              string `schema:"bg_color"`
	}

	GetFormHandler struct {
	}
)

func NewGetFormHandler() *GetFormHandler {
	return &GetFormHandler{}
}

func (h *GetFormHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetFormData) (templ.Component, error) {
	url := request.URL.Path
	templQuestionTypes := questionTypesForTempl(url, in.ActiveQuestionType)

	switch {
	case in.AnswerOptionFormType != "" && in.ActiveQuestionType != "" && in.BGColor != "":
		switch in.AnswerOptionFormType {
		case answerOptionFormTypeForm:
			return answerOptionForm(in.BGColor, url, in.ActiveQuestionType), nil
		case answerOptionFormTypeAdd:
			return answerOptionAdd(in.BGColor, url, in.ActiveQuestionType), nil
		}
	}

	switch in.ActiveQuestionType {
	case string(model.QuestionTypeMultipleChoice):
		return frontend_admin_question.QuestionForm(
			string(model.QuestionTypeMultipleChoice),
			submitUrl,
			questionTypes[model.QuestionTypeMultipleChoice].color,
			templQuestionTypes,
			[]templ.Component{
				answerOptionForm("indigo-500", url, in.ActiveQuestionType),
				answerOptionForm("pink-500", url, in.ActiveQuestionType),
				answerOptionAdd("blue-500", url, in.ActiveQuestionType),
				answerOptionAdd("red-500", url, in.ActiveQuestionType),
			},
		), nil
	case string(model.QuestionTypeChoice):
		return frontend_admin_question.QuestionForm(
			string(model.QuestionTypeChoice),
			submitUrl,
			questionTypes[model.QuestionTypeChoice].color,
			templQuestionTypes,
			[]templ.Component{
				answerOptionForm("orange-500", url, in.ActiveQuestionType),
				answerOptionForm("pink-500", url, in.ActiveQuestionType),
				answerOptionAdd("amber-500", url, in.ActiveQuestionType),
				answerOptionAdd("red-500", url, in.ActiveQuestionType),
			},
		), nil
	}

	return frontend.AdminPageComponent(
		formAddTitle,
		frontendComponents.Composition(
			frontendComponents.BackLink(listUrl),
			frontend_admin_question.QuestionForm(
				string(model.QuestionTypeChoice),
				submitUrl,
				questionTypes[model.QuestionTypeChoice].color,
				templQuestionTypes,
				[]templ.Component{
					answerOptionForm("orange-500", url, in.ActiveQuestionType),
					answerOptionForm("pink-500", url, in.ActiveQuestionType),
					answerOptionAdd("amber-500", url, in.ActiveQuestionType),
					answerOptionAdd("red-500", url, in.ActiveQuestionType),
				},
			),
		),
	), nil
}

func questionTypesForTempl(url string, activeQuestionType string) []frontend_admin_question.QuestionType {
	if activeQuestionType == "" {
		activeQuestionType = string(model.QuestionTypeChoice)
	}

	result := make([]frontend_admin_question.QuestionType, 0, len(questionTypes))
	for k, v := range questionTypes {
		code := string(k)
		item := frontend_admin_question.QuestionType{
			Name:      v.name,
			ActionUrl: fmt.Sprintf("%s?%s=%s", url, paramQuestionType, code),
			Color:     v.color,
		}
		if activeQuestionType == code {
			item.IsActive = true
		}

		result = append(result, item)
	}

	goSlices.SortFunc(result, func(i, j frontend_admin_question.QuestionType) int {
		return -1 * cmp.Compare(i.Name, j.Name)
	})

	return result
}

func answerOptionForm(color string, url string, activeQuestionType string) templ.Component {
	if activeQuestionType == "" {
		activeQuestionType = string(model.QuestionTypeChoice)
	}

	return frontend_admin_question.AnswerOptionForm(
		activeQuestionType,
		color,
		fmt.Sprintf(
			"%s?%s=%s&%s=%s&%s=%s",
			url,
			paramAnswerOptionType,
			answerOptionFormTypeAdd,
			paramBGColor,
			color,
			paramQuestionType,
			activeQuestionType,
		),
	)
}

func answerOptionAdd(color string, url string, activeQuestionType string) templ.Component {
	if activeQuestionType == "" {
		activeQuestionType = string(model.QuestionTypeChoice)
	}

	return frontend_admin_question.AnswerOptionAdd(
		color,
		fmt.Sprintf(
			"%s?%s=%s&%s=%s&%s=%s",
			url,
			paramAnswerOptionType,
			answerOptionFormTypeForm,
			paramBGColor,
			color,
			paramQuestionType,
			activeQuestionType,
		),
	)
}
