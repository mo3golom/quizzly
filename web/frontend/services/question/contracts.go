package question

import (
	"context"
	"github.com/a-h/templ"
	"github.com/google/uuid"
)

const (
	ListTypeFull    ListType = "full"    // full list with page wrapper
	ListTypeCompact ListType = "compact" // only list without page wrapper and header
)

type (
	ListType string

	ListOptions struct {
		Type            ListType
		SelectIsEnabled bool
	}

	Spec struct {
		QuestionIDs []uuid.UUID
		AuthorID    *uuid.UUID
	}

	Service interface {
		List(ctx context.Context, spec *Spec, options *ListOptions) (templ.Component, error)
	}
)
