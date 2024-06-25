package admin

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/pkg/structs"
	"quizzly/web/frontend/services/question"
	frontendIndex "quizzly/web/frontend/templ"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
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
	questionList, err := h.questionService.List(
		request.Context(),
		&question.Spec{
			AuthorID: structs.Pointer(uuid.New()),
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
		frontendAdminGame.Form(questionList),
	), nil
}
