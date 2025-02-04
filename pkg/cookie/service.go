package cookie

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"quizzly/pkg/variables"
	"time"
)

type (
	DefaultService struct {
		cookie *securecookie.SecureCookie
	}
)

func NewService(variablesRepo variables.Repository) Service {
	return &DefaultService{
		cookie: securecookie.New(
			[]byte(variablesRepo.GetString(variables.AuthSecretKey)),
			[]byte(variablesRepo.GetString(variables.AuthCookieBlockKey)),
		),
	}
}

func (s *DefaultService) Set(w http.ResponseWriter, key string, value string, ttl time.Duration) error {
	encoded, err := s.cookie.Encode(key, value)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     key,
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

func (s *DefaultService) Get(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", err
	}

	var value string
	err = s.cookie.Decode(key, cookie.Value, &value)
	if err != nil {
		return "", err
	}

	return value, err
}

func (s *DefaultService) Remove(w http.ResponseWriter, key string) {
	cookie := http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
