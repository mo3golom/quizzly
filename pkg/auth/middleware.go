package auth

import (
	"net/http"
)

type defaultMiddleware struct {
	repository           *defaultRepository
	encryptor            *defaultEncryptor
	forbiddenRedirectURL *string
}

func (s *defaultMiddleware) WithAuth(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromCookie(r)
		if token == "" {
			s.forbidden(w, r)
			return
		}

		encryptedToken, err := s.encryptor.Encrypt(string(token))
		if err != nil {
			s.forbidden(w, r)
			return
		}

		specificUser, err := s.repository.getUserByToken(
			r.Context(),
			Token(encryptedToken),
		)
		if err != nil {
			s.forbidden(w, r)
			return
		}

		authContext := DefaultAuthContext{
			Context: r.Context(),
			userID:  specificUser.id,
		}
		r = r.WithContext(authContext)

		delegate(w, r)
	}
}

func (s *defaultMiddleware) forbidden(w http.ResponseWriter, r *http.Request) {
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
