package login

import (
	"github.com/a-h/templ"
	"net/http"
	"quizzly/web/frontend/services/page"
	frontendLogin "quizzly/web/frontend/templ/admin/login"
)

const (
	loginPageTitle = "Login"
)

type (
	GetLoginPageHandler struct {
	}
)

func NewGetLoginPageHandler() *GetLoginPageHandler {
	return &GetLoginPageHandler{}
}

func (h *GetLoginPageHandler) Handle(_ http.ResponseWriter, request *http.Request, _ struct{}) (templ.Component, error) {
	return page.PublicIndexPage(
		request.Context(),
		loginPageTitle,
		frontendLogin.Page(
			frontendLogin.Form(""),
		),
	), nil
}
