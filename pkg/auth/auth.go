package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"quizzly/pkg/transactional"
	"time"
)

const (
	defaultTTL = 10 * time.Minute
)

type DefaultSimpleAuth struct {
	template   transactional.Template
	repository *defaultRepository
	generator  *defaultGenerator
	encryptor  *defaultEncryptor
	sender     Sender
}

func (a *DefaultSimpleAuth) SendLoginCode(ctx context.Context, email Email) error {
	code := a.generator.generateCode()

	encryptedData, err := a.encryptor.Encrypt(string(email))
	if err != nil {
		return err
	}
	encryptedEmail := Email(encryptedData)

	err = a.template.Execute(ctx, func(tx transactional.Tx) error {
		specificUser, err := a.repository.getUserByEmail(ctx, tx, encryptedEmail)
		if err != nil && !errors.Is(err, errUserNotFound) {
			return err
		}
		if errors.Is(err, errUserNotFound) {
			specificUser = &user{
				id:    uuid.New(),
				email: encryptedEmail,
			}
			err := a.repository.insertUser(ctx, tx, specificUser)
			if err != nil {
				return err
			}
		}

		err = a.repository.insertLoginCode(ctx, tx, &insertLoginCodeIn{
			userID:    specificUser.id,
			code:      code,
			expiresAt: time.Now().Add(defaultTTL),
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return a.sender.SendLoginCode(ctx, email, code)
}

func (a *DefaultSimpleAuth) Login(email Email, code LoginCode) (*uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}
