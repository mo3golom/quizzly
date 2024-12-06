package auth

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"net/http"
	"quizzly/pkg/transactional"
	"time"

	"github.com/google/uuid"
)

const (
	loginCodeDefaultTTL = 10 * time.Minute
	tokenDefaultTTL     = 336 * time.Hour
)

type DefaultSimpleAuth struct {
	template      transactional.Template
	repository    *defaultRepository
	generator     *defaultGenerator
	encryptor     *defaultEncryptor
	tokenService  *tokenService
	cookieService *cookieService
	sender        *sender
	cleaner       *DefaultCleaner
}

func NewSimpleAuth(
	db *sqlx.DB,
	template transactional.Template,
	config *Config,
) SimpleAuth {
	encryptor := &defaultEncryptor{
		secretKey: config.SecretKey,
	}
	repository := &defaultRepository{
		db: db,
	}

	return &DefaultSimpleAuth{
		template:   template,
		repository: repository,
		generator:  &defaultGenerator{},
		encryptor:  encryptor,
		sender: &sender{
			config: config,
		},
		tokenService:  newTokenService(config.SecretKey),
		cookieService: newCookieService(config.SecretKey, config.CookieBlockKey),
		cleaner: &DefaultCleaner{
			template:   template,
			repository: repository,
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

func (a *DefaultSimpleAuth) Login(
	ctx context.Context,
	w http.ResponseWriter,
	email Email,
	code LoginCode,
) error {
	encryptedData, err := a.encryptor.Encrypt(string(email))
	if err != nil {
		return err
	}
	encryptedEmail := Email(encryptedData)

	var token string
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

		token, err = a.tokenService.createToken(specificUser.id, tokenDefaultTTL)
		return err
	})
	if err != nil {
		return err
	}

	return a.cookieService.setToken(w, token, tokenDefaultTTL)
}

func (a *DefaultSimpleAuth) Logout(w http.ResponseWriter) {
	a.cookieService.removeToken(w)
}

func (a *DefaultSimpleAuth) Middleware(forbiddenRedirectURL ...string) SimpleAuthMiddleware {
	var url *string
	if len(forbiddenRedirectURL) > 0 {
		url = &forbiddenRedirectURL[0]
	}

	return &authMiddleware{
		repository:           a.repository,
		cookieService:        a.cookieService,
		tokenService:         a.tokenService,
		forbiddenRedirectURL: url,
	}
}

func (a *DefaultSimpleAuth) Cleaner() Cleaner {
	return a.cleaner
}
