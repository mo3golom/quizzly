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

const (
	cookiePlayerID   = "player-id"
	getPlayPageTitle = "Play"
)

type (
	GetPlayPageData struct {
		GameID  uuid.UUID `schema:"id"`
		Restart *bool     `schema:"restart"`
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
	game, err := h.gameUC.Get(request.Context(), in.GameID)
	if errors.Is(err, contracts.ErrGameNotFound) {
		return frontend.PublicPageComponent(
			getPlayPageTitle,
			frontendPublicGame.StartPage(),
		), nil
	}
	if err != nil {
		return nil, err
	}
	if game == nil || game.Status == model.GameStatusFinished {
		return frontend.PublicPageComponent(
			getPlayPageTitle,
			frontendPublicGame.StartPage(),
		), nil
	}

	playerID, err := h.getPlayerID(request)
	if err != nil {
		return nil, err
	}
	setPlayerID(writer, playerID)

	if in.Restart != nil && *in.Restart {
		err = h.sessionUC.Restart(request.Context(), game.ID, playerID)
		if err != nil {
			return nil, err
		}
	}

	session, err := h.sessionUC.GetCurrentState(request.Context(), in.GameID, playerID)
	if errors.Is(err, contracts.ErrQuestionQueueIsEmpty) {
		err = h.sessionUC.Finish(context.Background(), game.ID, playerID)
		if err != nil {
			return nil, err
		}

		return h.statistics(request.Context(), game, playerID)
	}
	if err != nil {
		return nil, err
	}

	if session.Status == model.SessionStatusFinished {
		return h.statistics(request.Context(), game, playerID)
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

	return frontend.PublicPageComponent(
		fmt.Sprintf("%s #%s", getPlayPageTitle, game.ID.String()),
		frontendPublicGame.Page(
			frontendPublicGame.QuestionComposition(
				game.ID,
				playerID,
				frontendComponents.GridLine(
					frontendPublicGame.Progress(&handlers.SessionProgress{
						Answered: int(session.Progress.Answered),
						Total:    int(session.Progress.Total),
					}),
					frontendPublicGame.Player(playerName),
				),
				frontendPublicGame.Question(&handlers.Question{
					ID:            session.CurrentQuestion.ID,
					Type:          session.CurrentQuestion.Type,
					Text:          session.CurrentQuestion.Text,
					ImageID:       session.CurrentQuestion.ImageID,
					AnswerOptions: answerOptions,
					Color:         frontend.ColorsMap[specificQuestionColor.Color][frontend.BgColor],
				}),
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

func (h *GetPlayPageHandler) statistics(ctx context.Context, game *model.Game, playerID uuid.UUID) (templ.Component, error) {
	stats, err := h.sessionUC.GetStatistics(ctx, game.ID, playerID)
	if err != nil {
		return nil, err
	}

	return frontend.PublicPageComponent(
		fmt.Sprintf("%s #%s", getPlayPageTitle, game.ID.String()),
		frontendComponents.CompositionMD(
			frontendPublicGame.ResultHeader(),
			frontendPublicGame.ResultStatistics(
				&handlers.SessionStatistics{
					QuestionsCount:      int(stats.QuestionsCount),
					CorrectAnswersCount: int(stats.CorrectAnswersCount),
				},
			),
			frontendPublicGame.ActionRestartGame(game.ID),
		),
	), nil
}
