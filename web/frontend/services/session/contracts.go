package session

import (
	"context"
	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type (
	ListOptions struct {
	}

	Spec struct {
		GameID uuid.UUID
	}

	Service interface {
		List(ctx context.Context, spec *Spec, options *ListOptions) ([]templ.Component, error)
	}
)
