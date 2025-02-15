package game

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/services/link"
	"quizzly/web/frontend/services/page"
	"quizzly/web/frontend/services/player"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

type (
	GetPlayPageData struct {
		GameID     *uuid.UUID `schema:"id"`
		CustomName *string    `schema:"name"`
	}

	GetPlayPageHandler struct {
		gameUC contracts.GameUsecase

		playerService player.Service

		service *service
	}
)

func NewGetPlayPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerService player.Service,
	linkService link.Service,
) *GetPlayPageHandler {
	return &GetPlayPageHandler{
		gameUC:        gameUC,
		playerService: playerService,
		service: &service{
			sessionUC:   sessionUC,
			linkService: linkService,
		},
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

	customPlayerName := ""
	if in.CustomName != nil {
		customPlayerName = *in.CustomName
	}
	currentPlayer, err := h.playerService.GetPlayer(writer, request, game.ID, customPlayerName)
	if err != nil {
		return nil, err
	}

	question, err := h.service.GetCurrentState(request.Context(), &getCurrentStateIn{
		game:       game,
		player:     currentPlayer,
		customName: in.CustomName,
	})
	if err != nil {
		return nil, err
	}

	return page.PublicIndexPage(
		request.Context(),
		gameTitle(game),
		frontendPublicGame.Page(question),
	), nil
}
