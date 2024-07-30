package auth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/jmoiron/sqlx"
	"quizzly/pkg/structs"
	"quizzly/pkg/transactional"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	loginCodeDefaultTTL = 10 * time.Minute
	tokenDefaultTTL     = 336 * time.Hour
)

type DefaultSimpleAuth struct {
	template   transactional.Template
	repository *defaultRepository
	generator  *defaultGenerator
	encryptor  *defaultEncryptor
	sender     Sender
	middleware SimpleAuthMiddleware
}

func NewSimpleAuth(
	db *sqlx.DB,
	template transactional.Template,
	encryptorConfig *EncryptorConfig,
	senderConfig *SenderConfig,
) SimpleAuth {
	encryptor := &defaultEncryptor{
		secretKey: encryptorConfig.SecretKey,
	}
	repository := &defaultRepository{
		db: db,
	}

	var sender Sender
	if senderConfig.Debug {
		sender = &DebugSender{}
	} else {
		sender = &DefaultSender{
			config: senderConfig,
		}
	}

	return &DefaultSimpleAuth{
		template:   template,
		repository: repository,
		generator:  &defaultGenerator{},
		encryptor:  encryptor,
		sender:     sender,
		middleware: &defaultMiddleware{
			repository: repository,
			encryptor:  encryptor,
		},
	}
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

		err = a.repository.upsertLoginCode(ctx, tx, &upsertLoginCodeIn{
			userID:    specificUser.id,
			code:      code,
			expiresAt: time.Now().Add(loginCodeDefaultTTL),
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

func (a *DefaultSimpleAuth) Login(ctx context.Context, email Email, code LoginCode) (*Token, error) {
	h := md5.New()
	h.Write([]byte(strings.ToLower(string(email))))
	token := Token(hex.EncodeToString(h.Sum(nil)))

	encryptedData, err := a.encryptor.Encrypt(string(email))
	if err != nil {
		return nil, err
	}
	encryptedEmail := Email(encryptedData)

	err = a.template.Execute(ctx, func(tx transactional.Tx) error {
		specificUser, err := a.repository.getUserByEmail(ctx, tx, encryptedEmail)
		if errors.Is(err, errUserNotFound) {
			return ErrLoginFailed
		}
		if err != nil {
			return err
		}

		specificCode, err := a.repository.getLoginCode(ctx, tx, &getLoginCodeIn{
			code:   code,
			userID: specificUser.id,
		})
		if errors.Is(err, errLoginCodeNotFound) {
			return ErrLoginFailed
		}
		if err != nil {
			return err
		}

		if specificUser.id != specificCode.userID {
			return ErrLoginFailed
		}

		encryptedToken, err := a.encryptor.Encrypt(string(token))
		if err != nil {
			return err
		}

		return a.repository.upsertToken(ctx, tx, &upsertTokenIn{
			token:     Token(encryptedToken),
			userID:    specificUser.id,
			expiresAt: time.Now().Add(tokenDefaultTTL),
		})
	})
	if err != nil {
		return nil, err
	}

	return structs.Pointer(token), nil
}

func (a *DefaultSimpleAuth) Middleware() SimpleAuthMiddleware {
	return a.middleware
}
