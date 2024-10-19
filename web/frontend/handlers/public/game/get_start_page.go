package game

import (
	"fmt"
	"github.com/a-h/templ"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/services/page"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

type (
	GetStartPageData struct {
		Warns []string `schema:"warn"`
	}

	GetStartPageHandler struct {
		uc contracts.GameUsecase
	}
)

func NewGetStartPageHandler(uc contracts.GameUsecase) *GetStartPageHandler {
	return &GetStartPageHandler{uc: uc}
}

func (h *GetStartPageHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetStartPageData) (templ.Component, error) {
	publicGames, err := h.uc.GetPublic(request.Context())
	if err != nil {
		return nil, err
	}

	return page.PublicIndexPage(
		request.Context(),
		"Играть в квизы",
		frontendPublicGame.Page(
			frontendPublicGame.StartPage(in.Warns...),
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
	), nil
}
