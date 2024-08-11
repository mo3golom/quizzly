package game

import (
	"github.com/a-h/templ"
	"net/http"
	frontendIndex "quizzly/web/frontend/templ"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
	frontendAdminQuestion "quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
)

const (
	title = "Новая игра"
)

type (
	GetFormHandler struct {
	}
)

func NewGetFormHandler() *GetFormHandler {
	return &GetFormHandler{}
}

func (h *GetFormHandler) Handle(_ http.ResponseWriter, _ *http.Request, _ struct{}) (templ.Component, error) {
	questionList := frontendAdminQuestion.QuestionListContainer(frontendAdminQuestion.ContainerOptions{
		WithSelect: true,
	})

	return frontendIndex.AdminPageComponent(
		title,
		frontendComponents.Composition(
			frontendComponents.BackLink(listUrl),
			frontendAdminGame.Form(questionList),
		),
	), nil
}
