package login

import (
	"errors"
	"github.com/a-h/templ"
	"net/http"
	"quizzly/pkg/auth"
	"quizzly/web/frontend/handlers"
	frontendLogin "quizzly/web/frontend/templ/admin/login"
	frontend_components "quizzly/web/frontend/templ/components"
	"time"
)

const (
	mainPageUrl = "/game/list"
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
		return nil, handlers.BadRequest(errors.New("email is required"))
	}

	if in.Email != nil && in.Code == nil {
		err := h.simpleAuth.SendLoginCode(request.Context(), *in.Email)
		if err != nil {
			return nil, err
		}

		return frontendLogin.Form(*in.Email), nil
	}

	token, err := h.simpleAuth.Login(request.Context(), *in.Email, *in.Code)
	if err != nil {
		return nil, handlers.BadRequest(err)
	}

	setToken(writer, *token)
	return frontend_components.Redirect(mainPageUrl), nil
}

func setToken(writer http.ResponseWriter, token auth.Token) {
	cookie := http.Cookie{
		Name:     auth.CookieToken,
		Value:    string(token),
		Path:     "/",
		Expires:  time.Now().Add(730001 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(writer, &cookie)
}
