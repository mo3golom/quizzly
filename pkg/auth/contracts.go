package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var (
	ErrLoginFailed = errors.New("login failed")
)

type (
	Email     string
	LoginCode int64

	Context interface {
		context.Context
		UserID() uuid.UUID
	}

	SimpleAuth interface {
		SendLoginCode(ctx context.Context, email Email) error
		Login(ctx context.Context, w http.ResponseWriter, email Email, code LoginCode) error
		Logout(w http.ResponseWriter)
		Middleware(forbiddenRedirectURL ...string) SimpleAuthMiddleware
		Cleaner() Cleaner
	}

	SimpleAuthMiddleware interface {
		Trace(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
		Auth(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
	}

	Cleaner interface {
		Name() string
		Perform(ctx context.Context) error
		DetermineInterval(ctx context.Context) (*time.Duration, error)
	}

	Config struct {
		SecretKey      string
		CookieBlockKey string // it should be 16 bytes (AES-128) or 32 bytes (AES-256) long
		Debug          bool

		FromEmail string
		Host      string
		Port      int64
		User      string
		Password  string
	}
)
