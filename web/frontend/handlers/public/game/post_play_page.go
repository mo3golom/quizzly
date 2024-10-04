package game

import (
	"context"
	"errors"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/services/player"
	frontend "quizzly/web/frontend/templ"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

type (
	PostPlayPageData struct {
		GameID     *uuid.UUID `schema:"id"`
		QuestionID uuid.UUID  `schema:"question-id"`
		PlayerID   uuid.UUID  `schema:"player-id"`
		Answers    []string   `schema:"answer"`
	}

	PostPlayPageHandler struct {
		gameUC    contracts.GameUsecase
		sessionUC contracts.SessionUsecase
		playerUC  contracts.PLayerUsecase

		playerService player.Service
	}
)

func NewPostPlayPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerUC contracts.PLayerUsecase,
	playerService player.Service,
) *PostPlayPageHandler {
	return &PostPlayPageHandler{
		gameUC:        gameUC,
		sessionUC:     sessionUC,
		playerUC:      playerUC,
		playerService: playerService,
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
		return frontendPublicGame.StartPage(), nil
	}

	playerID := in.PlayerID
	if playerID == uuid.Nil {
		currentPlayer, err := h.playerService.GetPlayer(writer, request)
		if err != nil {
			return nil, err
		}

		playerID = currentPlayer.ID
	}

	answerResult, err := h.sessionUC.AcceptAnswers(request.Context(), &contracts.AcceptAnswersIn{
		GameID:     game.ID,
		PlayerID:   playerID,
		QuestionID: in.QuestionID,
		Answers:    in.Answers,
	})
	if err != nil {
		return nil, err
	}

	answerComponent := buildAnswerComponent(answerResult, game.Settings.ShowRightAnswers)

	session, err := h.sessionUC.GetCurrentState(request.Context(), game.ID, playerID)
	if errors.Is(err, contracts.ErrQuestionQueueIsEmpty) {
		return h.finish(game, playerID, answerComponent)
	}
	if err != nil {
		return nil, err
	}

	specificQuestionColor := handlers.QuestionTypePublicColors
	answerOptions := make([]handlers.AnswerOption, 0, len(session.CurrentQuestion.AnswerOptions))
	for i, answerOption := range session.CurrentQuestion.AnswerOptions {
		answerOptions = append(answerOptions, handlers.AnswerOption{
			ID:    int64(answerOption.ID),
			Text:  answerOption.Answer,
			Color: frontend.ColorsMap[specificQuestionColor.AnswerOptionColors[i]][frontend.BgWithHoverColor],
		})
	}

	var playerName string
	players, err := h.playerUC.Get(request.Context(), []uuid.UUID{playerID})
	if err != nil {
		return nil, err
	}
	if len(players) > 0 {
		playerName = players[0].Name
	}

	var question templ.Component
	switch session.CurrentQuestion.Type {
	case model.QuestionTypeChoice, model.QuestionTypeOneOfChoice, model.QuestionTypeMultipleChoice:
		question = frontendPublicGame.Question(
			session.CurrentQuestion.ID,
			frontendPublicGame.QuestionBlock(session.CurrentQuestion.Text, session.CurrentQuestion.ImageID),
			frontendComponents.Composition(
				frontendPublicGame.AnswerChoiceDescription(session.CurrentQuestion.Type),
				frontendPublicGame.AnswerChoiceOptions(session.CurrentQuestion.Type, answerOptions),
			),
		)
	case model.QuestionTypeFillTheGap:
		question = frontendPublicGame.Question(
			session.CurrentQuestion.ID,
			frontendPublicGame.QuestionBlock(session.CurrentQuestion.Text, session.CurrentQuestion.ImageID),
			frontendPublicGame.AnswerTextInput(),
		)
	}

	return frontendPublicGame.QuestionForm(
		game.ID,
		playerID,
		frontendPublicGame.Header(game.Title),
		frontendComponents.GridLine(
			frontendPublicGame.Progress(&handlers.SessionProgress{
				Answered: int(session.Progress.Answered),
				Total:    int(session.Progress.Total),
			}),
			frontendPublicGame.Player(playerName),
		),
		question,
		answerComponent,
	), nil
}

func (h *PostPlayPageHandler) finish(
	game *model.Game,
	playerID uuid.UUID,
	answerComponent templ.Component,
) (templ.Component, error) {
	err := h.sessionUC.Finish(context.Background(), game.ID, playerID)
	if err != nil {
		return nil, err
	}

	return frontendComponents.Composition(
		frontendPublicGame.ResultLinkInput(resultsLink(game.ID, playerID)),
		answerComponent,
	), nil
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
