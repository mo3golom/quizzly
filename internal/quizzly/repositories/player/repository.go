package player

import (
	"context"
	"database/sql"
	"errors"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	sqlxPlayer struct {
		ID     uuid.UUID  `db:"id"`
		Name   string     `db:"name"`
		UserID *uuid.UUID `db:"user_id"`
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
		insert into player (id, user_id, name) values ($1, $2, $3) on conflict (id) do nothing
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.UserID, in.Name)
	return err
}

func (r *DefaultRepository) Get(ctx context.Context, id uuid.UUID) (*model.Player, error) {
	const query = ` 
		select id, user_id, name
		from player
		where id = $1
		limit 1
	`

	var result sqlxPlayer
	if err := r.db.GetContext(ctx, &result, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &model.Player{
		ID:     result.ID,
		Name:   result.Name,
		UserID: result.UserID,
	}, nil
}
