package game

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/services/link"
	"quizzly/web/frontend/services/page"
	"quizzly/web/frontend/services/player"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type (
	GetPlayPageData struct {
		GameID     *uuid.UUID `schema:"id"`
		CustomName *string    `schema:"name"`
	}

	GetPlayPageHandler struct {
		gameUC    contracts.GameUsecase
		sessionUC contracts.SessionUsecase
		playerUC  contracts.PLayerUsecase

		playerService player.Service

		linkService link.Service
	}
)

func NewGetPlayPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerUC contracts.PLayerUsecase,
	playerService player.Service,
	linkService link.Service,
) *GetPlayPageHandler {
	return &GetPlayPageHandler{
		gameUC:        gameUC,
		sessionUC:     sessionUC,
		playerUC:      playerUC,
		playerService: playerService,
		linkService:   linkService,
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
		return frontendComponents.Redirect("/?warn=Игра не найдена"), nil
	}
	if err != nil {
		return nil, err
	}
	if game.Status == model.GameStatusFinished {
		return frontendComponents.Redirect("/?warn=Игра уже завершена"), nil
	}
	if game.Status == model.GameStatusCreated {
		return frontendComponents.Redirect("/?warn=Игра еще не началась. Подождите немного или попросите автора запустить игру"), nil
	}

	customPlayerName := ""
	if in.CustomName != nil {
		customPlayerName = *in.CustomName
	}
	currentPlayer, err := h.playerService.GetPlayer(writer, request, game.ID, customPlayerName)
	if err != nil {
		return nil, err
	}

	if game.Settings.InputCustomName && !currentPlayer.NameUserEntered && in.CustomName == nil {
		return page.PublicIndexPage(
			request.Context(),
			h.getTitle(game),
			frontendPublicGame.Page(
				frontendPublicGame.NamePage(game.Title, game.ID),
			),
		), nil
	}

	session, err := h.sessionUC.GetCurrentState(request.Context(), *gameID, currentPlayer.ID)
	if errors.Is(err, contracts.ErrQuestionQueueIsEmpty) {
		err = h.sessionUC.Finish(context.Background(), game.ID, currentPlayer.ID)
		if err != nil {
			return nil, err
		}

		return frontendComponents.Redirect(h.linkService.GameResultsLink(game.ID, currentPlayer.ID)), nil
	}
	if err != nil {
		return nil, err
	}

	if session.Status == model.SessionStatusFinished {
		return frontendComponents.Redirect(h.linkService.GameResultsLink(game.ID, currentPlayer.ID)), nil
	}

	answerOptions := make([]handlers.AnswerOption, 0, len(session.CurrentQuestion.AnswerOptions))
	for _, answerOption := range session.CurrentQuestion.AnswerOptions {
		answerOptions = append(answerOptions, handlers.AnswerOption{
			ID:   int64(answerOption.ID),
			Text: answerOption.Answer,
		})
	}

	var playerName string
	players, err := h.playerUC.Get(request.Context(), []uuid.UUID{currentPlayer.ID})
	if err != nil {
		return nil, err
	}
	if len(players) > 0 {
		playerName = players[0].Name
	}

	return page.PublicIndexPage(
		request.Context(),
		h.getTitle(game),
		frontendPublicGame.Page(
			frontendPublicGame.QuestionForm(
				game.ID,
				currentPlayer.ID,
				frontendPublicGame.Header(game.Title),
				frontendComponents.GridLine(
					frontendPublicGame.Progress(&handlers.SessionProgress{
						Answered: int(session.Progress.Answered),
						Total:    int(session.Progress.Total),
					}),
					frontendPublicGame.Player(playerName),
				),
				frontendPublicGame.Question(
					session.CurrentQuestion.ID,
					frontendPublicGame.QuestionBlock(session.CurrentQuestion.Text, session.CurrentQuestion.ImageID),
					frontendComponents.Composition(
						frontendPublicGame.AnswerChoiceDescription(session.CurrentQuestion.Type),
						frontendPublicGame.AnswerChoiceOptions(session.CurrentQuestion.Type, answerOptions),
					),
				),
			),
		),
	), nil
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
