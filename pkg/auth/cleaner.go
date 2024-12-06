package auth

import (
	"context"
	"quizzly/pkg/structs"
	"quizzly/pkg/transactional"
	"time"
)

const (
	jobName     = "auth_login_codes_cleaner"
	jobInterval = time.Hour * 1
)

type DefaultCleaner struct {
	template   transactional.Template
	repository *defaultRepository
}

func (c *DefaultCleaner) Name() string {
	return jobName
}

func (c *DefaultCleaner) DetermineInterval(_ context.Context) (*time.Duration, error) {
	return structs.Pointer(jobInterval), nil
}

func (c *DefaultCleaner) Perform(ctx context.Context) error {
	return c.template.Execute(ctx, func(tx transactional.Tx) error {
		return c.repository.clearExpiredLoginCodes(ctx, tx)
	})
}
