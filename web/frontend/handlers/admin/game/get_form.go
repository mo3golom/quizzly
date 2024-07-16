package game

import (
	"github.com/a-h/templ"
	"net/http"
	"quizzly/pkg/auth"
	"quizzly/pkg/structs"
	"quizzly/web/frontend/services/question"
	frontendIndex "quizzly/web/frontend/templ"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
	frontendComponents "quizzly/web/frontend/templ/components"
)

const (
	title = "Новая игра"
)

type (
	GetFormHandler struct {
		questionService question.Service
	}
)

func NewGetFormHandler(questionService question.Service) *GetFormHandler {
	return &GetFormHandler{
		questionService: questionService,
	}
}

func (h *GetFormHandler) Handle(_ http.ResponseWriter, request *http.Request, _ struct{}) (templ.Component, error) {
	authContext := request.Context().(auth.Context)
	questionList, err := h.questionService.List(
		request.Context(),
		&question.Spec{
			AuthorID: structs.Pointer(authContext.UserID()),
		},
		&question.ListOptions{
			Type:            question.ListTypeCompact,
			SelectIsEnabled: true,
		},
	)
	if err != nil {
		return nil, err
	}

	return frontendIndex.AdminPageComponent(
		title,
		frontendComponents.Composition(
			frontendComponents.BackLink(listUrl),
			frontendAdminGame.Form(questionList),
		),
	), nil
}
