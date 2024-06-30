package login

import (
	"github.com/a-h/templ"
	"net/http"
	frontend "quizzly/web/frontend/templ"
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

func (h *GetLoginPageHandler) Handle(_ http.ResponseWriter, _ *http.Request, _ struct{}) (templ.Component, error) {
	return frontend.PublicPageComponent(
		loginPageTitle,
		frontendLogin.Page(
			frontendLogin.Form("", false),
		),
	), nil
}
