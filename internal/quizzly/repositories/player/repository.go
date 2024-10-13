package player

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	sqlxPlayer struct {
		ID              uuid.UUID  `db:"id"`
		Name            string     `db:"name"`
		NameUserEntered bool       `db:"name_user_entered"`
		UserID          *uuid.UUID `db:"user_id"`
	}

	DefaultRepository struct {
		db *sqlx.DB
	}
)

func NewRepository(db *sqlx.DB) Repository {
	return &DefaultRepository{
		db: db,
	}
}

func (r *DefaultRepository) Insert(ctx context.Context, tx transactional.Tx, in *model.Player) error {
	const query = ` 
		insert into player (id, user_id, name, name_user_entered) values ($1, $2, $3, $4) on conflict (id) do nothing
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.UserID, in.Name, in.NameUserEntered)
	return err
}

func (r *DefaultRepository) Update(ctx context.Context, tx transactional.Tx, in *model.Player) error {
	const query = ` 
		update player set 
		 user_id = $2, 
		 name = $3,
		 name_user_entered = $4
	    where id = $1
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.UserID, in.Name, in.NameUserEntered)
	return err
}

func (r *DefaultRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Player, error) {
	const query = ` 
		select id, user_id, name, name_user_entered
		from player
		where id = any($1)
	`

	var result []sqlxPlayer
	if err := r.db.SelectContext(ctx, &result, query, pq.Array(ids)); err != nil {
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
	if err := r.db.SelectContext(ctx, &result, query, pq.Array(userIDs)); err != nil {
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
