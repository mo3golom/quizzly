package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

const (
	CookieToken = "token"
)

var (
	ErrLoginFailed = errors.New("login failed")
)

type (
	Email     string
	LoginCode int64
	Token     string

	Context interface {
		context.Context
		UserID() uuid.UUID
	}

	SimpleAuth interface {
		SendLoginCode(ctx context.Context, email Email) error
		Login(ctx context.Context, email Email, code LoginCode) (*Token, error)
		Middleware() SimpleAuthMiddleware
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
		Debug     bool
	}

	EncryptorConfig struct {
		SecretKey string
	}
)
