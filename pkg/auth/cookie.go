package auth

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"time"
)

const (
	cookieJWT = "JWT"
)

type (
	cookieService struct {
		cookie *securecookie.SecureCookie
	}
)

func newCookieService(secretKey string, blockKey string) *cookieService {
	return &cookieService{
		cookie: securecookie.New([]byte(secretKey), []byte(blockKey)),
	}
}

func (s *cookieService) setToken(w http.ResponseWriter, token string, ttl time.Duration) error {
	encoded, err := s.cookie.Encode(cookieJWT, token)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     cookieJWT,
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(ttl),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	return nil
}

func (s *cookieService) getToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieJWT)
	if err != nil {
		return "", err
	}

	var value string
	err = s.cookie.Decode(cookieJWT, cookie.Value, &value)
	if err != nil {
		return "", err
	}

	return value, err
}

func (s *cookieService) removeToken(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     cookieJWT,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
