package auth

import (
	"context"
	"quizzly/pkg/transactional"
)

type DefaultCleaner struct {
	template   transactional.Template
	repository *defaultRepository
}

func (c *DefaultCleaner) ClearExpiredLoginCodes(ctx context.Context) error {
	return c.template.Execute(ctx, func(tx transactional.Tx) error {
		return c.repository.clearExpiredLoginCodes(ctx, tx)
	})
}
