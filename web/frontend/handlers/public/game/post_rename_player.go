package game

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/web/frontend/services/player"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

type (
	PostRenamePlayerData struct {
		Name string `schema:"name"`
	}

	PostRenamePlayerHandler struct {
		playerService player.Service
		playerUC      contracts.PLayerUsecase
	}
)

func NewPostRenamePlayerHandler(
	playerUC contracts.PLayerUsecase,
	playerService player.Service,
) *PostRenamePlayerHandler {
	return &PostRenamePlayerHandler{
		playerService: playerService,
		playerUC:      playerUC,
	}
}

func (h *PostRenamePlayerHandler) Handle(writer http.ResponseWriter, request *http.Request, in PostRenamePlayerData) (templ.Component, error) {
	var gameID *uuid.UUID
	var playerID *uuid.UUID
	if pathGameID := request.PathValue(pathValueGameID); pathGameID != "" {
		tempGameID, err := uuid.Parse(pathGameID)
		if err != nil {
			return nil, err
		}

		gameID = &tempGameID
	}
	if pathPlayerID := request.PathValue(pathValuePlayerID); pathPlayerID != "" {
		tempPLayerID, err := uuid.Parse(pathPlayerID)
		if err != nil {
			return nil, err
		}

		playerID = &tempPLayerID
	}

	if gameID == nil {
		return nil, errors.New("game id is required")
	}
	if playerID == nil {
		return nil, errors.New("player id is required")
	}

	currentPlayer, err := h.playerService.GetPlayer(writer, request, *gameID)
	if err != nil {
		return nil, err
	}

	if currentPlayer.ID != *playerID {
		return nil, errors.New("player not found")
	}

	if in.Name != "" {
		currentPlayer.Name = in.Name
		err = h.playerUC.Update(request.Context(), currentPlayer)
		if err != nil {
			return nil, err
		}
	}

	return frontendPublicGame.ActionRenamePlayer(
		*gameID,
		currentPlayer.ID,
		currentPlayer.Name,
	), nil
}
