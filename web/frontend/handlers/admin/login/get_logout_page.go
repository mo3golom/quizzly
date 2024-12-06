package login

import (
	"github.com/a-h/templ"
	"net/http"
	"quizzly/pkg/auth"
	frontend_components "quizzly/web/frontend/templ/components"
)

const (
	loginPageUrl = "/admin/login"
)

type (
	GetLogoutPageHandler struct {
		simpleAuth auth.SimpleAuth
	}
)

func NewGetLogoutPageHandler(
	simpleAuth auth.SimpleAuth,
) *GetLogoutPageHandler {
	return &GetLogoutPageHandler{
		simpleAuth: simpleAuth,
	}
}

func (h *GetLogoutPageHandler) Handle(writer http.ResponseWriter, _ *http.Request, _ struct{}) (templ.Component, error) {
	h.simpleAuth.Logout(writer)

	return frontend_components.Redirect(loginPageUrl), nil
}
