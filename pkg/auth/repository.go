package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"quizzly/pkg/transactional"
)

type (
	repository struct {
		db *sqlx.DB
	}
)

func (r *repository) insertUser(ctx context.Context, tx transactional.Tx, id uuid.UUID, email Email) error {
	const query = ` 
		insert into "user" (id, email) values ($1, $2) 
	    on conflict (email) do nothing
	`

	_, err := tx.ExecContext(ctx, query, id, email)
	return err
}

func (r *repository) insertLoginCode(ctx context.Context, tx transactional.Tx, in *insertLoginCodeIn) error {
	const query = ` 
		insert into user_auth_login_code (user_id, code, expires_at) values ($1, $2, $3) 
	    on conflict (user_id) do update set
			code=excluded.code,
			expires_at=excluded.expires_at
	`

	_, err := tx.ExecContext(ctx, query, in.userID, in.code, in.expiresAt)
	return err
}

func (r *repository) getLoginCode(ctx context.Context, tx transactional.Tx, in *getLoginCodeIn) (*loginCodeExtended, error) {
	const query = ` 
		select code from user_auth_login_code
		where user_id=$1 and code=$2 and expires_at > now()
		limit 1
	`

	var result struct {
		UserID uuid.UUID `db:"user_id"`
		Code   LoginCode `db:"code"`
	}
	err := tx.GetContext(ctx, &result, query, in.userID, in.code)
	if err != nil {
		return nil, err
	}

	return &loginCodeExtended{
		userID: result.UserID,
		code:   result.Code,
	}, nil
}
