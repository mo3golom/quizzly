package supabase

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type (
	Auth interface {
		OTP(email string) error
		LoginOTP(w http.ResponseWriter, token string, redirectTo string) error
		Logout(w http.ResponseWriter, r *http.Request) error

		MiddlewareTrace(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
		MiddlewareAuth(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
	}

	AuthContext interface {
		context.Context
		UserID() uuid.UUID
	}
)
