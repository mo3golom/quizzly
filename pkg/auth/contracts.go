package auth

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type (
	Email     string
	LoginCode int64

	SimpleAuth interface {
		SendLoginCode(ctx context.Context, email Email) error
		Login(ctx context.Context, email Email, code LoginCode) (*uuid.UUID, error)
		ClearLoginCodes(ctx context.Context) error
	}

	SimpleAuthMiddleware interface {
		WithAuth(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
	}
)
