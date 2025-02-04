package login

import (
	"github.com/google/uuid"
	"net/http"
	"quizzly/pkg/supabase"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/services/page"
	frontend_components "quizzly/web/frontend/templ/components"
	frontend_admin_login "quizzly/web/frontend/templ/public/login"

	"github.com/a-h/templ"
)

const (
	loginPageTitle = "Login"
	mainPageUrl    = "/admin/game/list"
)

type (
	GetLoginPageData struct {
		Email *string `schema:"email"`

		Token      *string `schema:"token"`
		RedirectTo string  `schema:"redirect_to"`
	}

	GetLoginPageHandler struct {
		authClient supabase.Auth
	}
)

func NewGetLoginPageHandler(authClient supabase.Auth) *GetLoginPageHandler {
	return &GetLoginPageHandler{
		authClient: authClient,
	}
}

func (h *GetLoginPageHandler) Handle(writer http.ResponseWriter, request *http.Request, in GetLoginPageData) (templ.Component, error) {
	authContext, ok := request.Context().(supabase.AuthContext)
	if ok && authContext.UserID() != uuid.Nil {
		return frontend_components.Redirect(mainPageUrl), nil
	}

	if in.Token != nil {
		err := h.authClient.LoginOTP(writer, *in.Token, in.RedirectTo)
		if err != nil {
			return nil, err
		}

		return frontend_components.Redirect(mainPageUrl), nil
	}

	if in.Email != nil {
		err := h.authClient.OTP(*in.Email)
		if err != nil {
			return nil, handlers.BadRequest(err)
		}

		return frontend_admin_login.Form(*in.Email), nil
	}

	return page.PublicIndexPage(
		request.Context(),
		loginPageTitle,
		frontend_admin_login.Page(
			frontend_admin_login.Form(""),
		),
	), nil
}
