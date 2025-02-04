package game

import (
	"fmt"
	"github.com/a-h/templ"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/supabase"
	frontend_components "quizzly/web/frontend/templ/components"
)

type (
	GetCreateHandler struct {
		uc contracts.GameUsecase
	}
)

func NewGetCreateHandler(uc contracts.GameUsecase) *GetCreateHandler {
	return &GetCreateHandler{
		uc: uc,
	}
}

func (h *GetCreateHandler) Handle(_ http.ResponseWriter, request *http.Request, _ struct{}) (templ.Component, error) {
	authContext := request.Context().(supabase.AuthContext)
	gameID, err := h.uc.Create(
		request.Context(),
		&contracts.CreateGameIn{
			AuthorID: authContext.UserID(),
			Type:     model.GameTypeAsync,
		},
	)
	if err != nil {
		return nil, err
	}

	return frontend_components.Redirect(fmt.Sprintf("/admin/game/%s", gameID.String())), nil
}
