package auth

import (
	"context"
	"errors"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
)

const (
	loginCodeDefaultTTL = 10 * time.Minute
	tokenDefaultTTL     = 336 * time.Hour
)

type DefaultSimpleAuth struct {
	trm           trm.Manager
	repository    *defaultRepository
	generator     *defaultGenerator
	encryptor     *defaultEncryptor
	tokenService  *tokenService
	cookieService *cookieService
	sender        *sender
	cleaner       *DefaultCleaner
}

func NewSimpleAuth(
	sqlx *sqlx.DB,
	trm trm.Manager,
	config *Config,
) SimpleAuth {
	trmsqlxGetter := trmsqlx.DefaultCtxGetter

	encryptor := &defaultEncryptor{
		secretKey: config.SecretKey,
	}
	repository := &defaultRepository{
		sqlx: sqlx,
		tx:   trmsqlxGetter,
	}

	return &DefaultSimpleAuth{
		trm:        trm,
		repository: repository,
		generator:  &defaultGenerator{},
		encryptor:  encryptor,
		sender: &sender{
			config: config,
		},
		tokenService:  newTokenService(config.SecretKey),
		cookieService: newCookieService(config.SecretKey, config.CookieBlockKey),
		cleaner: &DefaultCleaner{
			trm:        trm,
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

	err = a.trm.Do(ctx, func(ctx context.Context) error {
		specificUser, err := a.repository.getUserByEmail(ctx, encryptedEmail)
		if err != nil && !errors.Is(err, errUserNotFound) {
			return err
		}
		if errors.Is(err, errUserNotFound) {
			specificUser = &user{
				id:    uuid.New(),
				email: encryptedEmail,
			}
			err := a.repository.insertUser(ctx, specificUser)
			if err != nil {
				return err
			}
		}

		err = a.repository.upsertLoginCode(ctx, &upsertLoginCodeIn{
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
	err = a.trm.Do(ctx, func(ctx context.Context) error {
		specificUser, err := a.repository.getUserByEmail(ctx, encryptedEmail)
		if errors.Is(err, errUserNotFound) {
			return ErrLoginFailed
		}
		if err != nil {
			return err
		}

		specificCode, err := a.repository.getLoginCode(ctx, &getLoginCodeIn{
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
