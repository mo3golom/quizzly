package game

import (
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/services/link"
	"quizzly/web/frontend/services/player"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type (
	PostPlayPageData struct {
		GameID     *uuid.UUID `schema:"id"`
		QuestionID uuid.UUID  `schema:"question-id"`
		Answers    []string   `schema:"answer"`
	}

	PostPlayPageHandler struct {
		gameUC        contracts.GameUsecase
		sessionUC     contracts.SessionUsecase
		playerService player.Service

		service *service
	}
)

func NewPostPlayPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerService player.Service,
	linkService link.Service,
) *PostPlayPageHandler {
	return &PostPlayPageHandler{
		gameUC:        gameUC,
		sessionUC:     sessionUC,
		playerService: playerService,
		service: &service{
			sessionUC:   sessionUC,
			linkService: linkService,
		},
	}
}

func (h *PostPlayPageHandler) Handle(writer http.ResponseWriter, request *http.Request, in PostPlayPageData) (templ.Component, error) {
	gameID := in.GameID
	if pathGameID := request.PathValue(pathValueGameID); pathGameID != "" {
		tempGameID, err := uuid.Parse(pathGameID)
		if err != nil {
			return nil, err
		}

		gameID = &tempGameID
	}

	game, err := h.gameUC.Get(request.Context(), *gameID)
	if err != nil {
		return nil, err
	}
	if game.Status == model.GameStatusFinished {
		return frontendComponents.Redirect("/?warn=Игра уже завершена"), nil
	}

	currentPlayer, err := h.playerService.GetPlayer(writer, request, game.ID)
	if err != nil {
		return nil, err
	}

	answerResult, err := h.sessionUC.AcceptAnswers(request.Context(), &contracts.AcceptAnswersIn{
		GameID:     game.ID,
		PlayerID:   currentPlayer.ID,
		QuestionID: in.QuestionID,
		Answers:    in.Answers,
	})
	if err != nil {
		return nil, err
	}

	answerComponent := buildAnswerComponent(answerResult, game.Settings.ShowRightAnswers)

	return h.service.GetCurrentState(
		request.Context(),
		&getCurrentStateIn{
			game:   game,
			player: currentPlayer,
		},
		answerComponent,
	)
}

func buildAnswerComponent(answerResult *contracts.AcceptAnswersOut, displayRightAnswers bool) templ.Component {
	var rightAnswers []string
	if displayRightAnswers {
		rightAnswers = slices.SafeMap(answerResult.RightAnswers, func(in model.AnswerOption) string {
			return in.Answer
		})
	}

	return frontendPublicGame.Answer(answerResult.IsCorrect, rightAnswers...)
}
