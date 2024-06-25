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
	}

	SimpleAuthMiddleware interface {
		WithAuth(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
	}

	Sender interface {
		SendLoginCode(ctx context.Context, to Email, code LoginCode) error
	}

	SenderConfig struct {
		FromEmail Email
		Host      string
		Port      int64
		User      string
		Password  string
	}
)
