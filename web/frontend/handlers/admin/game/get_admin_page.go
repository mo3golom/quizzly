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
	listUrl      = "/admin/game/list"
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
	gameID := in.GameID
	if pathGameID := request.PathValue(pathValueGameID); pathGameID != "" {
		tempGameID, err := uuid.Parse(pathGameID)
		if err != nil {
			return nil, err
		}

		gameID = &tempGameID
	}

	if gameID == nil {
		return frontend.AdminPageComponent(
			getPageTitle,
			frontendAdminGame.NotFound(),
		), nil
	}

	game, err := h.uc.Get(request.Context(), *gameID)
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

	return frontend.AdminPageComponent(
		getPageTitle,
		frontendAdminGame.Page(
			frontendComponents.BackLink(listUrl),
			frontendAdminGame.Header(
				&handlers.Game{
					ID:        game.ID,
					Status:    game.Status,
					Title:     game.Title,
					CreatedAt: game.CreatedAt,
				},
			),
			frontendAdminGame.Settings(&handlers.GameSettings{
				ShuffleQuestions: game.Settings.ShuffleQuestions,
				ShuffleAnswers:   game.Settings.ShuffleAnswers,
				ShowRightAnswers: game.Settings.ShowRightAnswers,
			}),
			frontendAdminGame.Invite(gameLink(game.ID, request)),
			frontendAdminGame.Statistics(
				&handlers.GameStatistics{
					QuestionsCount:    int(stats.QuestionsCount),
					ParticipantsCount: int(stats.ParticipantsCount),
					CompletionRate:    int(stats.CompletionRate),
				},
			),
			frontendComponents.Tabs(
				uuid.New(),
				frontendComponents.Tab{
					Name:    "Вопросы",
					Content: questionList,
				},
				frontendComponents.Tab{
					Name:    "Участники",
					Content: frontendAdminGame.SessionListContainer(game.ID),
				},
			),
		),
	), nil
}
