package game

import (
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/pkg/supabase"
	frontend "quizzly/web/frontend/templ"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
	frontendComponents "quizzly/web/frontend/templ/components"
	"sort"

	"github.com/a-h/templ"
)

const (
	getListTitle = "Список игр"
)

type (
	GetListHandler struct {
		uc contracts.GameUsecase
	}
)

func NewGetListHandler(uc contracts.GameUsecase) *GetListHandler {
	return &GetListHandler{
		uc: uc,
	}
}

func (h *GetListHandler) Handle(_ http.ResponseWriter, request *http.Request, _ struct{}) (templ.Component, error) {
	authContext := request.Context().(supabase.AuthContext)
	games, err := h.uc.GetByAuthor(request.Context(), authContext.UserID())
	if err != nil {
		return nil, err
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].CreatedAt.After(games[j].CreatedAt)
	})

	components := make([]templ.Component, 0, len(games)+1)
	components = append(components, frontendComponents.Header(
		getListTitle,
		frontendAdminGame.ActionAddNewGame(),
	))
	for _, game := range games {
		game := game
		components = append(components, frontendAdminGame.GameListItem(convertModelGameToHandlersGame(&game)))
	}

	return frontend.AdminPageComponent(
		getListTitle,
		frontendComponents.Composition(
			components...,
		),
	), nil
}
