package jobs

import (
	"context"
	"time"
)

type (
	ContextEnricher interface {
		Enrich(context.Context) context.Context
	}

	Runner interface {
		Run(context.Context) error
		Register(Job) error
		RegisterAll(...Job) error
	}

	Job interface {
		Name() string
		Perform(ctx context.Context) error
		DetermineInterval(ctx context.Context) (*time.Duration, error)
	}
)
