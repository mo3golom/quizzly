package game

import (
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/services/page"
	"quizzly/web/frontend/services/player"
	frontend "quizzly/web/frontend/templ"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

const (
	getPlayResultsTitle            = "Результат игры"
	getPlayResultsShareDescription = "Мой результат %d из %d в игре %s. Сыграй и ты!"
)

type (
	GetPlayResultsPageData struct {
		GameID   *uuid.UUID `schema:"game_id"`
		PlayerID *uuid.UUID `schema:"id"`
	}

	GetPlayResultsPageHandler struct {
		gameUC    contracts.GameUsecase
		sessionUC contracts.SessionUsecase
		playerUC  contracts.PLayerUsecase

		playerService player.Service
	}
)

func NewGetPlayResultsPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerUC contracts.PLayerUsecase,
	playerService player.Service,
) *GetPlayResultsPageHandler {
	return &GetPlayResultsPageHandler{
		gameUC:        gameUC,
		sessionUC:     sessionUC,
		playerUC:      playerUC,
		playerService: playerService,
	}
}

func (h *GetPlayResultsPageHandler) Handle(writer http.ResponseWriter, request *http.Request, in GetPlayResultsPageData) (templ.Component, error) {
	gameID := in.GameID
	if pathGameID := request.PathValue(pathValueGameID); pathGameID != "" {
		tempGameID, err := uuid.Parse(pathGameID)
		if err != nil {
			return nil, err
		}

		gameID = &tempGameID
	}

	playerID := in.PlayerID
	if pathPlayerID := request.PathValue(pathValuePlayerID); pathPlayerID != "" {
		tempPLayerID, err := uuid.Parse(pathPlayerID)
		if err != nil {
			return nil, err
		}

		playerID = &tempPLayerID
	}

	game, err := h.gameUC.Get(request.Context(), *gameID)
	if err != nil {
		return nil, err
	}

	currentPlayer, err := h.playerService.GetPlayer(writer, request, game.ID)
	if err != nil {
		return nil, err
	}

	stats, err := h.sessionUC.GetStatistics(request.Context(), game.ID, *playerID)
	if errors.Is(err, contracts.ErrSessionNotFinished) {
		return page.PublicIndexPage(
			request.Context(),
			h.getTitle(game.Title),
			frontendComponents.StatusMessage("Ой ей... Игра еще не завершена"),
		), nil
	}
	if err != nil {
		return nil, err
	}

	var playerName string
	players, err := h.playerUC.Get(request.Context(), []uuid.UUID{*playerID})
	if err != nil {
		return nil, err
	}
	if len(players) > 0 {
		playerName = players[0].Name
	}

	resultPlayer := frontendPublicGame.ResultPlayer(playerName)
	actions := make([]templ.Component, 0, 2)
	if currentPlayer.ID == *playerID {
		actions = append(actions, frontendPublicGame.ActionRestartGame(game.ID))
		actions = append(actions, frontendPublicGame.ActionShareResult(h.getShareTitle(game.Title, stats.CorrectAnswersCount, stats.QuestionsCount)))

		resultPlayer = frontendPublicGame.ResultPlayer(
			playerName,
			frontendPublicGame.ActionRenamePlayer(game.ID, currentPlayer.ID, playerName),
		)
	} else {
		actions = append(actions, frontendPublicGame.ActionPlayGame(game.ID))
	}

	publicGames, err := h.gameUC.GetPublic(request.Context())
	if err != nil {
		return nil, err
	}

	return page.PublicIndexPage(
		request.Context(),
		h.getTitle(game.Title),
		frontendPublicGame.Page(
			frontendPublicGame.ResultHeader(game.Title),
			resultPlayer,
			frontendPublicGame.ResultStatistics(
				&handlers.SessionStatistics{
					QuestionsCount:      int(stats.QuestionsCount),
					CorrectAnswersCount: int(stats.CorrectAnswersCount),
				},
			),
			frontendComponents.GridLine(actions...),
			frontendPublicGame.ResultAdditional(
				frontendComponents.DividerVerticalLight("Или",
					frontendPublicGame.PublicGameComposition(
						slices.SafeMap(publicGames[:5], func(game model.Game) templ.Component {
							title := fmt.Sprintf("Игра от %s", game.CreatedAt.Format("02.01.2006"))
							if game.Title != nil {
								title = *game.Title
							}

							return frontendPublicGame.PublicGame(title, game.ID)
						})...,
					),
					frontendPublicGame.CreateGame(),
				),
			),
		),
		frontend.OpenGraph{
			Title: h.getShareTitle(game.Title, stats.CorrectAnswersCount, stats.QuestionsCount),
			URL:   resultsLink(game.ID, *playerID, request),
		}), nil
}

func (h *GetPlayResultsPageHandler) getTitle(gameTitle *string) string {
	if gameTitle == nil {
		return getPlayResultsTitle
	}

	return fmt.Sprintf(`%s "%s"`, getPlayResultsTitle, *gameTitle)
}

func (h *GetPlayResultsPageHandler) getShareTitle(gameTitle *string, correctAnswersCount int64, answersCount int64) string {
	title := ""
	if gameTitle != nil {
		title = fmt.Sprintf(`"%s"`, *gameTitle)
	}

	return fmt.Sprintf(getPlayResultsShareDescription, correctAnswersCount, answersCount, title)
}
