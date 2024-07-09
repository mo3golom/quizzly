package game

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/services/question"
	"quizzly/web/frontend/services/session"
	frontend "quizzly/web/frontend/templ"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
	frontendComponents "quizzly/web/frontend/templ/components"
)

const (
	getPageTitle = "Управление игрой"
)

type (
	GetAdminPageData struct {
		GameID *uuid.UUID `schema:"id"`
	}

	GetAdminPageHandler struct {
		uc              contracts.GameUsecase
		questionService question.Service
		sessionService  session.Service
	}
)

func NewGetPageHandler(
	uc contracts.GameUsecase,
	questionService question.Service,
	sessionService session.Service,
) *GetAdminPageHandler {
	return &GetAdminPageHandler{
		uc:              uc,
		questionService: questionService,
		sessionService:  sessionService,
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

	questionList, err := h.questionService.List(
		request.Context(),
		&question.Spec{
			QuestionIDs: questionIDs,
		},
		&question.ListOptions{
			Type:            question.ListTypeCompact,
			SelectIsEnabled: false,
		},
	)
	if err != nil {
		return nil, err
	}

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
			frontendAdminGame.Header(
				&handlers.Game{
					ID:        game.ID,
					Status:    game.Status,
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
						[]string{"Имя", "Процент прохождения", "Статус прохождения"},
						sessionList...,
					),
				},
			),
		),
	), nil
}
