package auth

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"net/http"
)

type authMiddleware struct {
	forbiddenRedirectURL *string
	repository           *defaultRepository
	encryptor            *defaultEncryptor
}

func (s *authMiddleware) WithEnrich(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enrichContextFn := func(r *http.Request) context.Context {
			token := getTokenFromCookie(r)
			if token == "" {
				return r.Context()
			}

			encryptedToken, err := s.encryptor.Encrypt(string(token))
			if err != nil {
				return r.Context()
			}

			specificUser, err := s.repository.getUserByToken(
				r.Context(),
				Token(encryptedToken),
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

func (s *authMiddleware) WithAuth(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return s.WithEnrich(func(w http.ResponseWriter, r *http.Request) {
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

func getTokenFromCookie(request *http.Request) Token {
	cookie, err := request.Cookie(CookieToken)
	if err != nil {
		return ""
	}

	return Token(cookie.Value)
}
