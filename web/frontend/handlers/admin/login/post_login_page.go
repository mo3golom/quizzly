package login

import (
	"errors"
	"github.com/a-h/templ"
	"net/http"
	"quizzly/pkg/auth"
	"quizzly/web/frontend/handlers"
	frontendLogin "quizzly/web/frontend/templ/admin/login"
	frontend_components "quizzly/web/frontend/templ/components"
)

const (
	mainPageUrl = "/admin/game/list"
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

	err := h.simpleAuth.Login(request.Context(), writer, *in.Email, *in.Code)
	if err != nil {
		return nil, handlers.BadRequest(err)
	}

	return frontend_components.Redirect(mainPageUrl), nil
}
