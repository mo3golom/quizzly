package cookie

import (
	"net/http"
	"time"
)

type (
	Service interface {
		Set(w http.ResponseWriter, key string, value string, ttl time.Duration) error
		Get(r *http.Request, key string) (string, error)
		Remove(w http.ResponseWriter, key string)
	}
)
