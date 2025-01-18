package auth

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type authMiddleware struct {
	forbiddenRedirectURL *string
	repository           *defaultRepository
	tokenService         *tokenService
	cookieService        *cookieService
}

func (s *authMiddleware) Trace(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enrichContextFn := func(r *http.Request) context.Context {
			token, err := s.cookieService.getToken(r)
			if err != nil || token == "" {
				return r.Context()
			}

			err = s.tokenService.verifyToken(token)
			if err != nil {
				return r.Context()
			}

			userID, err := s.tokenService.getUserID(token)
			if err != nil {
				return r.Context()
			}

			specificUser, err := s.repository.getUserByID(
				r.Context(),
				userID,
			)
			if err != nil {
				return r.Context()
			}

			return DefaultAuthContext{
				Context: r.Context(),
				userID:  specificUser.id,
			}
		}

		r = r.WithContext(enrichContextFn(r))
		delegate(w, r)
	}
}

func (s *authMiddleware) Auth(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return s.Trace(func(w http.ResponseWriter, r *http.Request) {
		authCtx, ok := r.Context().(Context)
		if !ok || authCtx.UserID() == uuid.Nil {
			s.forbidden(w, r)
			return
		}

		delegate(w, r)
	})
}

func (s *authMiddleware) forbidden(w http.ResponseWriter, r *http.Request) {
	if s.forbiddenRedirectURL != nil {
		http.Redirect(w, r, *s.forbiddenRedirectURL, http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusForbidden)
}
