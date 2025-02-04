package login

import (
	"net/http"
	"quizzly/pkg/supabase"
	frontend_components "quizzly/web/frontend/templ/components"

	"github.com/a-h/templ"
)

const (
	loginPageUrl = "/login"
)

type (
	GetLogoutPageHandler struct {
		authClient supabase.Auth
	}
)

func NewGetLogoutPageHandler(
	authClient supabase.Auth,
) *GetLogoutPageHandler {
	return &GetLogoutPageHandler{
		authClient: authClient,
	}
}

func (h *GetLogoutPageHandler) Handle(writer http.ResponseWriter, request *http.Request, _ struct{}) (templ.Component, error) {
	err := h.authClient.Logout(writer, request)
	if err != nil {
		return nil, err
	}

	return frontend_components.Redirect(loginPageUrl), nil
}
