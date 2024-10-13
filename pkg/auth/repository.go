package auth

import (
	"context"
	"database/sql"
	"errors"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	errUserNotFound      = errors.New("user not found")
	errLoginCodeNotFound = errors.New("login code not found")
)

type (
	sqlxUser struct {
		ID    uuid.UUID `db:"id"`
		Email Email     `db:"email"`
	}

	defaultRepository struct {
		db *sqlx.DB
	}
)

func (r *defaultRepository) insertUser(ctx context.Context, tx transactional.Tx, in *user) error {
	const query = ` 
		insert into "user" (id, email) values ($1, $2) 
	    on conflict (email) do nothing
	`

	_, err := tx.ExecContext(ctx, query, in.id, in.email)
	return err
}

func (r *defaultRepository) getUserByEmail(ctx context.Context, tx transactional.Tx, email Email) (*user, error) {
	const query = ` 
		select id, email from "user"
		where email = $1
		limit 1
	`

	var result sqlxUser
	err := tx.GetContext(ctx, &result, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user{
		id:    result.ID,
		email: result.Email,
	}, nil
}

func (r *defaultRepository) upsertLoginCode(ctx context.Context, tx transactional.Tx, in *upsertLoginCodeIn) error {
	const query = ` 
		insert into user_auth_login_code (user_id, code, expires_at) values ($1, $2, $3) 
	    on conflict (user_id) do update set
			code=excluded.code,
			expires_at=excluded.expires_at,
			updated_at = now()
	`

	_, err := tx.ExecContext(ctx, query, in.userID, in.code, in.expiresAt)
	return err
}

func (r *defaultRepository) getLoginCode(ctx context.Context, tx transactional.Tx, in *getLoginCodeIn) (*loginCodeExtended, error) {
	const query = ` 
		select user_id, code from user_auth_login_code
		where user_id=$1 and code=$2 and expires_at > now()
		limit 1
	`

	var result struct {
		UserID uuid.UUID `db:"user_id"`
		Code   LoginCode `db:"code"`
	}
	err := tx.GetContext(ctx, &result, query, in.userID, in.code)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errLoginCodeNotFound
	}
	if err != nil {
		return nil, err
	}

	return &loginCodeExtended{
		userID: result.UserID,
		code:   result.Code,
	}, nil
}

func (r *defaultRepository) clearExpiredLoginCodes(ctx context.Context, tx transactional.Tx) error {
	const query = ` 
		delete from user_auth_login_code
		where expires_at < now()
	`

	_, err := tx.ExecContext(ctx, query)
	return err
}

func (r *defaultRepository) upsertToken(ctx context.Context, tx transactional.Tx, in *upsertTokenIn) error {
	const query = ` 
		insert into user_auth_token (user_id, token, expires_at) values ($1, $2, $3) 
	    on conflict (user_id) do update set
			token = excluded.token,
			expires_at = excluded.expires_at,
			updated_at = now()
	`

	_, err := tx.ExecContext(ctx, query, in.userID, in.token, in.expiresAt)
	return err
}

func (r *defaultRepository) getUserByToken(ctx context.Context, token Token) (*user, error) {
	const query = ` 
		select u.id, u.email from "user" u
		join user_auth_token uat on uat.user_id = u.id
		where uat.token = $1 and uat.expires_at > now()
		limit 1
	`

	var result sqlxUser
	err := r.db.GetContext(ctx, &result, query, token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user{
		id:    result.ID,
		email: result.Email,
	}, nil
}
