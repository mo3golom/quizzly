package game

import (
	"context"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

const (
	defaultLimit int64 = 10_000_000
)

type (
	DefaultRepository struct {
		sqlx *sqlx.DB
		tx   *trmsqlx.CtxGetter
	}
)

func NewRepository(sqlx *sqlx.DB, tx *trmsqlx.CtxGetter) Repository {
	return &DefaultRepository{sqlx: sqlx, tx: tx}
}

func (r *DefaultRepository) db(ctx context.Context) trmsqlx.Tr {
	return r.tx.DefaultTrOrDB(ctx, r.sqlx)
}
