package auth

import (
	"context"
	"database/sql"
	"errors"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
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
		sqlx *sqlx.DB
		tx   *trmsqlx.CtxGetter
	}
)

func (r *defaultRepository) db(ctx context.Context) trmsqlx.Tr {
	return r.tx.DefaultTrOrDB(ctx, r.sqlx)
}

func (r *defaultRepository) getUserByEmail(ctx context.Context, email Email) (*user, error) {
	const query = ` 
		select id, email from "user"
		where email = $1
		limit 1
	`

	var result sqlxUser
	err := r.db(ctx).GetContext(ctx, &result, query, email)
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

func (r *defaultRepository) getUserByID(ctx context.Context, id uuid.UUID) (*user, error) {
	const query = ` 
		select u.id, u.email from "user" u
		where u.id = $1
		limit 1
	`

	var result sqlxUser
	err := r.db(ctx).GetContext(ctx, &result, query, id)
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

func (r *defaultRepository) insertUser(ctx context.Context, in *user) error {
	const query = ` 
		insert into "user" (id, email) values ($1, $2) 
	    on conflict (email) do nothing
	`

	_, err := r.db(ctx).ExecContext(ctx, query, in.id, in.email)
	return err
}

func (r *defaultRepository) upsertLoginCode(ctx context.Context, in *upsertLoginCodeIn) error {
	const query = ` 
		insert into user_auth_login_code (user_id, code, expires_at) values ($1, $2, $3) 
	    on conflict (user_id) do update set
			code=excluded.code,
			expires_at=excluded.expires_at,
			updated_at = now()
	`

	_, err := r.db(ctx).ExecContext(ctx, query, in.userID, in.code, in.expiresAt)
	return err
}

func (r *defaultRepository) getLoginCode(ctx context.Context, in *getLoginCodeIn) (*loginCodeExtended, error) {
	const query = ` 
		select user_id, code from user_auth_login_code
		where user_id=$1 and code=$2 and expires_at > now()
		limit 1
	`

	var result struct {
		UserID uuid.UUID `db:"user_id"`
		Code   LoginCode `db:"code"`
	}
	err := r.db(ctx).GetContext(ctx, &result, query, in.userID, in.code)
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

func (r *defaultRepository) clearExpiredLoginCodes(ctx context.Context) error {
	const query = ` 
		delete from user_auth_login_code
		where expires_at < now()
	`

	_, err := r.db(ctx).ExecContext(ctx, query)
	return err
}
