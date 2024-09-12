package session

import (
	"context"
	"github.com/a-h/templ"
	"github.com/google/uuid"
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
		List(ctx context.Context, spec *Spec, page int64, limit int64) (*ListOut, error)
	}
)
