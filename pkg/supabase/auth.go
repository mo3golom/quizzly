package supabase

import (
	"context"
	"fmt"
	"net/http"
	"quizzly/pkg/cookie"
	"quizzly/pkg/variables"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/auth-go/types"
)

const (
	cookieJWT = "JWT"
)

type DefaultAuth struct {
	client auth.Client
	cookie cookie.Service
}

func NewAuth(
	cookieService cookie.Service,
	variablesRepo variables.Repository,
) Auth {
	client := auth.New(
		variablesRepo.GetString(variables.SupabaseReference),
		variablesRepo.GetString(variables.SupabaseAnonKey),
	)

	return &DefaultAuth{
		client: client,
		cookie: cookieService,
	}
}

func (a *DefaultAuth) OTP(email string) error {
	return a.client.OTP(types.OTPRequest{
		Email:      email,
		CreateUser: true,
	})
}

func (a *DefaultAuth) LoginOTP(w http.ResponseWriter, token string, redirectTo string) error {
	result, err := a.client.Verify(types.VerifyRequest{
		Type:       types.VerificationTypeMagiclink,
		Token:      token,
		RedirectTo: redirectTo,
	})
	if err != nil {
		return err
	}

	if result.Error != "" {
		return fmt.Errorf("%s (%s): %s", result.Error, result.ErrorCode, result.ErrorDescription)
	}

	return a.cookie.Set(w, cookieJWT, result.AccessToken, time.Duration(result.ExpiresIn)*time.Second)
}

func (a *DefaultAuth) Logout(w http.ResponseWriter, r *http.Request) error {
	token, err := a.cookie.Get(r, cookieJWT)
	if err != nil {
		return err
	}

	err = a.client.WithToken(token).Logout()
	if err != nil {
		return err
	}

	a.cookie.Remove(w, cookieJWT)
	return nil
}

func (a *DefaultAuth) MiddlewareTrace(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enrichContextFn := func(r *http.Request) context.Context {
			token, err := a.cookie.Get(r, cookieJWT)
			if err != nil {
				return r.Context()
			}

			user, err := a.client.WithToken(token).GetUser()
			if err != nil {
				return r.Context()
			}

			return DefaultAuthContext{
				Context: r.Context(),
				userID:  user.ID,
			}
		}

		r = r.WithContext(enrichContextFn(r))
		delegate(w, r)
	}
}

func (a *DefaultAuth) MiddlewareAuth(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return a.MiddlewareTrace(func(w http.ResponseWriter, r *http.Request) {
		authCtx, ok := r.Context().(AuthContext)
		if !ok || authCtx.UserID() == uuid.Nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		delegate(w, r)
	})
}
