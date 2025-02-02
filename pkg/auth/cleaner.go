package auth

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"quizzly/pkg/structs"
	"time"
)

const (
	jobName     = "auth_login_codes_cleaner"
	jobInterval = time.Hour * 1
)

type DefaultCleaner struct {
	trm        trm.Manager
	repository *defaultRepository
}

func (c *DefaultCleaner) Name() string {
	return jobName
}

func (c *DefaultCleaner) DetermineInterval(_ context.Context) (*time.Duration, error) {
	return structs.Pointer(jobInterval), nil
}

func (c *DefaultCleaner) Perform(ctx context.Context) error {
	return c.trm.Do(ctx, func(ctx context.Context) error {
		return c.repository.clearExpiredLoginCodes(ctx)
	})
}
