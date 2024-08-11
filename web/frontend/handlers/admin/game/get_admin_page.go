package game

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/services/session"
	frontend "quizzly/web/frontend/templ"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
	frontendAdminQuestion "quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
)

const (
	getPageTitle = "Управление игрой"
	listUrl      = "/game/list"
)

type (
	GetAdminPageData struct {
		GameID *uuid.UUID `schema:"id"`
	}

	GetAdminPageHandler struct {
		uc             contracts.GameUsecase
		sessionService session.Service
	}
)

func NewGetPageHandler(
	uc contracts.GameUsecase,
	sessionService session.Service,
) *GetAdminPageHandler {
	return &GetAdminPageHandler{
		uc:             uc,
		sessionService: sessionService,
	}
}

func (h *GetAdminPageHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetAdminPageData) (templ.Component, error) {
	if in.GameID == nil {
		return frontend.AdminPageComponent(
			getPageTitle,
			frontendAdminGame.NotFound(),
		), nil
	}

	game, err := h.uc.Get(request.Context(), *in.GameID)
	if err != nil {
		return nil, err
	}

	stats, err := h.uc.GetStatistics(request.Context(), game.ID)
	if err != nil {
		return nil, err
	}

	questionIDs, err := h.uc.GetQuestions(request.Context(), game.ID)
	if err != nil {
		return nil, err
	}

	questionList := frontendAdminQuestion.QuestionListContainer(frontendAdminQuestion.ContainerOptions{
		QuestionIDs: questionIDs,
	})

	sessionList, err := h.sessionService.List(
		request.Context(),
		&session.Spec{
			GameID: game.ID,
		},
		&session.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	return frontend.AdminPageComponent(
		getPageTitle,
		frontendAdminGame.Page(
			frontendComponents.BackLink(listUrl),
			frontendAdminGame.Header(
				&handlers.Game{
					ID:        game.ID,
					Status:    game.Status,
					Title:     game.Title,
					Link:      getGameLink(game.ID, request),
					CreatedAt: game.CreatedAt,
				},
			),
			frontendAdminGame.Statistics(
				getGameLink(game.ID, request),
				&handlers.GameStatistics{
					QuestionsCount:    int(stats.QuestionsCount),
					ParticipantsCount: int(stats.ParticipantsCount),
					CompletionRate:    int(stats.CompletionRate),
				},
			),
			frontendAdminGame.Settings(&handlers.GameSettings{
				ShuffleQuestions: game.Settings.ShuffleQuestions,
				ShuffleAnswers:   game.Settings.ShuffleAnswers,
			}),
			frontendComponents.Tabs(
				uuid.New(),
				frontendComponents.Tab{
					Name:    "Вопросы",
					Content: questionList,
				},
				frontendComponents.Tab{
					Name: "Участники",
					Content: frontendComponents.Table(
						[]string{
							"Имя",
							"Процент прохождения",
							"Дата старта",
							"Дата последнего ответа",
							"Статус прохождения"},
						sessionList...,
					),
				},
			),
		),
	), nil
}
