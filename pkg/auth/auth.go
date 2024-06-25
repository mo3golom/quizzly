package auth

import (
	"context"
	"github.com/google/uuid"
	"quizzly/pkg/transactional"
)

type DefaultSimpleAuth struct {
	template transactional.Template
	repo     *repository
}

func (a *DefaultSimpleAuth) SendLoginCode(ctx context.Context, email Email) error {
	err := a.template.Execute(ctx, func(tx transactional.Tx) error {
		userID := uuid.New()
		err := a.repo.insertUser(ctx, tx, userId, email)
		if err != nil {
			return err
		}

		err := a.repo.insertLoginCode(ctx, tx, &insertLoginCodeIn{
			userID:
		})
	})
}

func (a *DefaultSimpleAuth) Login(email Email, code LoginCode) (*uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (a *DefaultSimpleAuth) ClearLoginCodes() error {
	//TODO implement me
	panic("implement me")
}
