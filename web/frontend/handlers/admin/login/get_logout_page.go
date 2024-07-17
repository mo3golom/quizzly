package login

import (
	"github.com/a-h/templ"
	"net/http"
	"quizzly/pkg/auth"
	frontend_components "quizzly/web/frontend/templ/components"
	"time"
)

const (
	loginPageUrl = "/login"
)

type (
	GetLogoutPageHandler struct{}
)

func NewGetLogoutPageHandler() *GetLogoutPageHandler {
	return &GetLogoutPageHandler{}
}

func (h *GetLogoutPageHandler) Handle(writer http.ResponseWriter, _ *http.Request, _ struct{}) (templ.Component, error) {
	removeToken(writer)

	return frontend_components.Redirect(loginPageUrl), nil
}

func removeToken(writer http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     auth.CookieToken,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(writer, &cookie)
}
