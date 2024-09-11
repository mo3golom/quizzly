package game

import (
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/handlers"
	frontend "quizzly/web/frontend/templ"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

const (
	getPlayResultsTitle            = "Результат игры"
	getPlayResultsShareDescription = "Я сыграл в игру %s, и вот мой результат!"
)

type (
	GetPlayResultsPageData struct {
		GameID   uuid.UUID `schema:"game_id"`
		PlayerID uuid.UUID `schema:"id"`
	}

	GetPlayResultsPageHandler struct {
		gameUC    contracts.GameUsecase
		sessionUC contracts.SessionUsecase
		playerUC  contracts.PLayerUsecase
	}
)

func NewGetPlayResultsPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerUC contracts.PLayerUsecase,
) *GetPlayResultsPageHandler {
	return &GetPlayResultsPageHandler{
		gameUC:    gameUC,
		sessionUC: sessionUC,
		playerUC:  playerUC,
	}
}

func (h *GetPlayResultsPageHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetPlayResultsPageData) (templ.Component, error) {
	playerID := getPlayerID(request)

	game, err := h.gameUC.Get(request.Context(), in.GameID)
	if err != nil {
		return nil, err
	}

	stats, err := h.sessionUC.GetStatistics(request.Context(), in.GameID, in.PlayerID)
	if errors.Is(err, contracts.ErrSessionNotFinished) {
		return frontend.PublicPageComponent(
			h.getTitle(game.Title),
			frontendComponents.StatusMessage("Ой ей... Игра еще не завершена"),
		), nil
	}
	if err != nil {
		return nil, err
	}

	actions := make([]templ.Component, 0, 2)
	if playerID == in.PlayerID {
		actions = append(actions, frontendPublicGame.ActionRestartGame(in.GameID))
		actions = append(actions, frontendPublicGame.ActionShareResult(getResultsLink(in.GameID, in.PlayerID, request), h.getShareTitle(game.Title)))
	} else {
		actions = append(actions, frontendPublicGame.ActionPlayGame(in.GameID))
	}

	var playerName string
	players, err := h.playerUC.Get(request.Context(), []uuid.UUID{in.PlayerID})
	if err != nil {
		return nil, err
	}
	if len(players) > 0 {
		playerName = players[0].Name
	}

	return frontend.PublicPageComponent(
		h.getTitle(game.Title),
		frontendPublicGame.Page(
			frontendPublicGame.ResultHeader(game.Title),
			frontendPublicGame.ResultPlayer(playerName),
			frontendPublicGame.ResultStatistics(
				&handlers.SessionStatistics{
					QuestionsCount:      int(stats.QuestionsCount),
					CorrectAnswersCount: int(stats.CorrectAnswersCount),
				},
			),
			frontendComponents.GridLine(actions...),
		),
		frontend.OpenGraph{
			Title: h.getShareTitle(game.Title),
			URL:   getResultsLink(in.GameID, in.PlayerID, request),
		}), nil
}

func (h *GetPlayResultsPageHandler) getTitle(gameTitle *string) string {
	if gameTitle == nil {
		return getPlayResultsTitle
	}

	return fmt.Sprintf(`%s "%s"`, getPlayResultsTitle, *gameTitle)
}

func (h *GetPlayResultsPageHandler) getShareTitle(gameTitle *string) string {
	if gameTitle != nil {
		return fmt.Sprintf(
			getPlayResultsShareDescription,
			fmt.Sprintf(`"%s"`, *gameTitle),
		)
	}

	return fmt.Sprintf(getPlayResultsShareDescription, "")
}