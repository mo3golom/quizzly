package game

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/handlers"
	frontend "quizzly/web/frontend/templ"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
	"time"
)

type (
	GetPlayPageData struct {
		GameID *uuid.UUID `schema:"id"`
	}

	GetPlayPageHandler struct {
		gameUC    contracts.GameUsecase
		sessionUC contracts.SessionUsecase
		playerUC  contracts.PLayerUsecase
	}
)

func NewGetPlayPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerUC contracts.PLayerUsecase,
) *GetPlayPageHandler {
	return &GetPlayPageHandler{
		gameUC:    gameUC,
		sessionUC: sessionUC,
		playerUC:  playerUC,
	}
}

func (h *GetPlayPageHandler) Handle(writer http.ResponseWriter, request *http.Request, in GetPlayPageData) (templ.Component, error) {
	gameID := in.GameID
	if pathGameID := request.PathValue(pathValueGameID); pathGameID != "" {
		tempGameID, err := uuid.Parse(pathGameID)
		if err != nil {
			return nil, err
		}

		gameID = &tempGameID
	}

	game, err := h.gameUC.Get(request.Context(), *gameID)
	if errors.Is(err, contracts.ErrGameNotFound) {
		return frontend.PublicPageComponent(
			h.getTitle(game),
			frontendPublicGame.StartPage("Игра не найдена."),
		), nil
	}
	if err != nil {
		return nil, err
	}
	if game.Status == model.GameStatusFinished {
		return frontend.PublicPageComponent(
			h.getTitle(game),
			frontendPublicGame.StartPage("Игра уже завершена."),
		), nil
	}
	if game.Status == model.GameStatusCreated {
		return frontend.PublicPageComponent(
			h.getTitle(game),
			frontendPublicGame.StartPage("Игра еще не началась. Подождите немного или попросите автора запустить игру."),
		), nil
	}

	playerID, err := h.getPlayerID(request)
	if err != nil {
		return nil, err
	}
	setPlayerID(writer, playerID)

	session, err := h.sessionUC.GetCurrentState(request.Context(), *gameID, playerID)
	if errors.Is(err, contracts.ErrQuestionQueueIsEmpty) {
		err = h.sessionUC.Finish(context.Background(), game.ID, playerID)
		if err != nil {
			return nil, err
		}

		return frontendComponents.Redirect(resultsLink(game.ID, playerID)), nil
	}
	if err != nil {
		return nil, err
	}

	if session.Status == model.SessionStatusFinished {
		return frontendComponents.Redirect(resultsLink(game.ID, playerID)), nil
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

	return frontend.PublicPageComponent(
		h.getTitle(game),
		frontendPublicGame.Page(
			frontendPublicGame.QuestionForm(
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
			),
		),
	), nil
}

func (h *GetPlayPageHandler) getPlayerID(request *http.Request) (uuid.UUID, error) {
	playerID := getPlayerID(request)
	if playerID == uuid.Nil {
		newPlayerID := uuid.New()
		err := h.playerUC.Create(request.Context(), &model.Player{
			ID:   newPlayerID,
			Name: namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate(),
		})
		if err != nil {
			return uuid.Nil, err
		}

		playerID = newPlayerID
	}

	return playerID, nil
}

func (h *GetPlayPageHandler) getTitle(game *model.Game) string {
	if game == nil {
		return "Игра не найдена"
	}

	title := fmt.Sprintf("Игра от %s", game.CreatedAt.Format("02.01.2006"))
	if game.Title != nil {
		title = fmt.Sprintf(`Игра "%s"`, *game.Title)
	}

	return title
}
