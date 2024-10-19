package page

import (
	"context"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"quizzly/pkg/auth"
	frontend "quizzly/web/frontend/templ"
)

func PublicIndexPage(ctx context.Context, title string, body templ.Component, openGraph ...frontend.OpenGraph) templ.Component {
	showAdminLink := false
	if authCtx, ok := ctx.(auth.Context); ok && authCtx.UserID() != uuid.Nil {
		showAdminLink = true
	}

	return frontend.PublicPageComponent(
		title,
		body,
		showAdminLink,
		openGraph...,
	)
}