package game

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/services/link"
	frontend "quizzly/web/frontend/templ"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
)

const (
	getPageTitle = "Управление игрой"
)

type (
	setting struct {
		slug  string
		text  string
		hint  string
		value func(settings *model.GameSettings) bool
	}

	GetGamePageData struct {
		GameID *uuid.UUID `schema:"id"`
	}

	GetGamePageHandler struct {
		service *service
	}
)

func NewGetPageHandler(
	uc contracts.GameUsecase,
	linkService link.Service,
) *GetGamePageHandler {
	return &GetGamePageHandler{
		service: &service{
			uc:          uc,
			linkService: linkService,
		},
	}
}

func (h *GetGamePageHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetGamePageData) (templ.Component, error) {
	gameID := in.GameID
	if pathGameID := request.PathValue(pathValueGameID); pathGameID != "" {
		tempGameID, err := uuid.Parse(pathGameID)
		if err != nil {
			return nil, err
		}

		gameID = &tempGameID
	}

	if gameID == nil {
		return frontend.AdminPageComponent(
			getPageTitle,
			frontendAdminGame.NotFound(),
		), nil
	}

	result, err := h.service.getGamePage(request, *gameID)
	if err != nil {
		return nil, err
	}

	return frontend.AdminPageComponent(
		getPageTitle,
		result,
	), nil
}
