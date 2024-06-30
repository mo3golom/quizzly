package login

import (
	"github.com/a-h/templ"
	"net/http"
	"quizzly/pkg/auth"
	frontendLogin "quizzly/web/frontend/templ/admin/login"
	"time"
)

type (
	PostLoginPageData struct {
		Email *auth.Email     `schema:"email"`
		Code  *auth.LoginCode `schema:"code"`
	}

	PostLoginPageHandler struct {
		simpleAuth auth.SimpleAuth
	}
)

func NewPostLoginPageHandler(
	simpleAuth auth.SimpleAuth,
) *PostLoginPageHandler {
	return &PostLoginPageHandler{
		simpleAuth: simpleAuth,
	}
}

func (h *PostLoginPageHandler) Handle(writer http.ResponseWriter, request *http.Request, in PostLoginPageData) (templ.Component, error) {
	if in.Email == nil {
		return frontendLogin.Form("", false), nil
	}

	if in.Email != nil && in.Code == nil {
		err := h.simpleAuth.SendLoginCode(request.Context(), *in.Email)
		if err != nil {
			return frontendLogin.Form(*in.Email, false, err), nil
		}

		return frontendLogin.Form(*in.Email, false), nil
	}

	token, err := h.simpleAuth.Login(request.Context(), *in.Email, *in.Code)
	if err != nil {
		return frontendLogin.Form(*in.Email, false, err), nil
	}

	setToken(writer, *token)
	return frontendLogin.Form(*in.Email, true), nil
}

func setToken(writer http.ResponseWriter, token auth.Token) {
	cookie := http.Cookie{
		Name:     auth.CookieToken,
		Value:    string(token),
		Path:     "/",
		Expires:  time.Now().Add(336 * time.Hour),
		MaxAge:   1209600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(writer, &cookie)
}
