package player

import (
	"context"
	"database/sql"
	"errors"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
)

type (
	sqlxPlayer struct {
		ID              uuid.UUID  `db:"id"`
		Name            string     `db:"name"`
		NameUserEntered bool       `db:"name_user_entered"`
		UserID          *uuid.UUID `db:"user_id"`
	}

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

func (r *DefaultRepository) Insert(ctx context.Context, in *model.Player) error {
	const query = ` 
		insert into player (id, user_id, name, name_user_entered) values ($1, $2, $3, $4) on conflict (id) do nothing
	`

	_, err := r.db(ctx).ExecContext(ctx, query, in.ID, in.UserID, in.Name, in.NameUserEntered)
	return err
}

func (r *DefaultRepository) Update(ctx context.Context, in *model.Player) error {
	const query = ` 
		update player set 
		 user_id = $2, 
		 name = $3,
		 name_user_entered = $4
	    where id = $1
	`

	_, err := r.db(ctx).ExecContext(ctx, query, in.ID, in.UserID, in.Name, in.NameUserEntered)
	return err
}

func (r *DefaultRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Player, error) {
	const query = ` 
		select id, user_id, name, name_user_entered
		from player
		where id = any($1)
	`

	var result []sqlxPlayer
	if err := r.db(ctx).SelectContext(ctx, &result, query, pq.Array(ids)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return slices.SafeMap(result, func(in sqlxPlayer) model.Player {
		return model.Player{
			ID:              in.ID,
			Name:            in.Name,
			NameUserEntered: in.NameUserEntered,
			UserID:          in.UserID,
		}
	}), nil
}

func (r *DefaultRepository) GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]model.Player, error) {
	const query = ` 
		select id, user_id, name, name_user_entered
		from player
		where user_id = any($1)
	`

	var result []sqlxPlayer
	if err := r.db(ctx).SelectContext(ctx, &result, query, pq.Array(userIDs)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return slices.SafeMap(result, func(in sqlxPlayer) model.Player {
		return model.Player{
			ID:              in.ID,
			Name:            in.Name,
			NameUserEntered: in.NameUserEntered,
			UserID:          in.UserID,
		}
	}), nil
}
