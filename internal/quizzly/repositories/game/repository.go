package game

import (
	"github.com/jmoiron/sqlx"
)

const (
	defaultLimit int64 = 10_000_000
)

type (
	DefaultRepository struct {
		db *sqlx.DB
	}
)

func NewRepository(db *sqlx.DB) Repository {
	return &DefaultRepository{db: db}
}
