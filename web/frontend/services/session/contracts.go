package session

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
)

type (
	ListOut struct {
		Result     []templ.Component
		TotalCount int64
	}

	Spec struct {
		GameID uuid.UUID
	}

	Service interface {
		List(request *http.Request, spec *Spec, page int64, limit int64) (*ListOut, error)
	}
)
