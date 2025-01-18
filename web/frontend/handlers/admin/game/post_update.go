package game

import (
	"errors"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"

	frontendComponents "quizzly/web/frontend/templ/components"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type (
	PostUpdateData struct {
		ShuffleQuestions *bool   `schema:"shuffle_questions"`
		ShuffleAnswers   *bool   `schema:"shuffle_answers"`
		ShowRightAnswers *bool   `schema:"show_right_answers"`
		InputCustomName  *bool   `schema:"input_custom_name"`
		IsPrivate        *bool   `schema:"is_private"`
		Title            *string `schema:"title"`
	}

	PostUpdateHandler struct {
		uc contracts.GameUsecase
	}
)

func NewPostUpdateHandler(uc contracts.GameUsecase) *PostUpdateHandler {
	return &PostUpdateHandler{uc: uc}
}

func (h *PostUpdateHandler) Handle(_ http.ResponseWriter, request *http.Request, in PostUpdateData) (templ.Component, error) {
	pathGameID := request.PathValue(pathValueGameID)
	if pathGameID == "" {
		return nil, errors.New("game_id is absent in path")
	}

	gameID, err := uuid.Parse(pathGameID)
	if err != nil {
		return nil, err
	}

	game, err := h.uc.Get(request.Context(), gameID)
	if err != nil {
		return nil, err
	}

	fill(&in, game)
	err = h.uc.Update(request.Context(), game)
	if err != nil {
		return nil, err
	}

	return frontendComponents.Composition(), nil
}

func fill(in *PostUpdateData, specificGame *model.Game) {
	if in.Title != nil {
		specificGame.Title = in.Title
	}
	if in.ShuffleQuestions != nil {
		specificGame.Settings.ShuffleQuestions = *in.ShuffleQuestions
	}
	if in.ShuffleAnswers != nil {
		specificGame.Settings.ShuffleAnswers = *in.ShuffleAnswers
	}
	if in.ShowRightAnswers != nil {
		specificGame.Settings.ShowRightAnswers = *in.ShowRightAnswers
	}
	if in.InputCustomName != nil {
		specificGame.Settings.InputCustomName = *in.InputCustomName
	}
	if in.IsPrivate != nil {
		specificGame.Settings.IsPrivate = *in.IsPrivate
	}
}
