package auth

import (
	"net/http"
)

type defaultMiddleware struct {
	repository *defaultRepository
	encryptor  *defaultEncryptor
}

func (s *defaultMiddleware) WithAuth(delegate func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromCookie(r)
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		encryptedToken, err := s.encryptor.Encrypt(string(token))
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		specificUser, err := s.repository.getUserByToken(
			r.Context(),
			Token(encryptedToken),
		)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
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

func getTokenFromCookie(request *http.Request) Token {
	cookie, err := request.Cookie(CookieToken)
	if err != nil {
		return ""
	}

	return Token(cookie.Value)
}
