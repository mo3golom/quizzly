package game

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/auth"
	"quizzly/pkg/structs/collections/slices"
	frontend_components "quizzly/web/frontend/templ/components"
)

type (
	PostCreateData struct {
		Questions        []string `schema:"question"`
		ShuffleQuestions bool     `schema:"shuffle_questions"`
		ShuffleAnswers   bool     `schema:"shuffle_answers"`
		ShowRightAnswers bool     `schema:"show_right_answers"`
		InputCustomName  bool     `schema:"input_custom_name"`
		IsPrivate        bool     `schema:"is_private"`
		Title            *string  `schema:"title"`
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
	if len(in.Questions) == 0 {
		return nil, fmt.Errorf("questions are required")
	}

	authContext := request.Context().(auth.Context)
	gameID, err := h.uc.Create(
		request.Context(),
		&contracts.CreateGameIn{
			AuthorID: authContext.UserID(),
			Type:     model.GameTypeAsync,
			Title:    in.Title,
			Settings: model.GameSettings{
				IsPrivate:        in.IsPrivate,
				ShuffleQuestions: in.ShuffleQuestions,
				ShuffleAnswers:   in.ShuffleAnswers,
				ShowRightAnswers: in.ShowRightAnswers,
				InputCustomName:  in.InputCustomName,
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

	return frontend_components.Redirect(fmt.Sprintf("/admin/game/%s", gameID.String())), nil
}
