package admin

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/handlers"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
)

type (
	PostCreateData struct {
		Questions        []string `schema:"question"`
		ShuffleQuestions bool     `schema:"shuffle_questions"`
		ShuffleAnswers   bool     `schema:"shuffle_answers"`
	}

	PostCreateHandler struct {
		uc contracts.GameUsecase
	}
)

func NewPostCreateHandler(uc contracts.GameUsecase) *PostCreateHandler {
	return &PostCreateHandler{
		uc: uc,
	}
}

func (h *PostCreateHandler) Handle(_ http.ResponseWriter, request *http.Request, in PostCreateData) (templ.Component, error) {
	gameID, err := h.uc.Create(
		request.Context(),
		&contracts.CreateGameIn{
			Type: model.GameTypeAsync,
			Settings: model.GameSettings{
				IsPrivate:        false,
				ShuffleQuestions: in.ShuffleQuestions,
				ShuffleAnswers:   in.ShuffleAnswers,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	questionIDs, err := slices.Map(in.Questions, func(id string) (uuid.UUID, error) {
		return uuid.Parse(id)
	})
	if err != nil {
		return nil, err
	}

	err = h.uc.AddQuestion(request.Context(), gameID, questionIDs...)
	if err != nil {
		return nil, err
	}

	return frontendAdminGame.Page(
		frontendAdminGame.Header(
			&handlers.Game{
				ID:     gameID,
				Status: model.GameStatusCreated,
				Link:   getGameLink(gameID, request),
			},
		),
	), nil
}
