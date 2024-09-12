package game

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/web/frontend/services/session"
	frontendComponents "quizzly/web/frontend/templ/components"
)

const (
	defaultLimit = 50
)

type (
	GetSessionListData struct {
		GameID     uuid.UUID `schema:"game_id"`
		PageNumber *int64    `schema:"page_number"`
	}

	GetSessionListHandler struct {
		sessionService session.Service
	}
)

func NewGetSessionListHandler(sessionService session.Service) *GetSessionListHandler {
	return &GetSessionListHandler{
		sessionService: sessionService,
	}
}

func (h *GetSessionListHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetSessionListData) (templ.Component, error) {
	page := int64(1)
	if in.PageNumber != nil {
		page = *in.PageNumber
	}

	sessionList, err := h.sessionService.List(
		request.Context(),
		&session.Spec{
			GameID: in.GameID,
		},
		page,
		defaultLimit,
	)
	if err != nil {
		return nil, err
	}

	return frontendComponents.Composition(
		frontendComponents.Table(
			[]string{
				"Имя",
				"Процент прохождения",
				"Дата старта",
				"Дата последнего ответа",
				"Статус прохождения"},
			sessionList.Result...,
		),
		frontendComponents.Pagination(
			page,
			sessionList.TotalCount,
			defaultLimit,
		),
	), nil
}
