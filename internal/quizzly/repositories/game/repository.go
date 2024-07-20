package game

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"
	"time"

	"github.com/google/uuid"
)

type (
	sqlxGame struct {
		ID                       uuid.UUID `db:"id"`
		AuthorID                 uuid.UUID `db:"author_id"`
		Status                   string    `db:"status"`
		Type                     string    `db:"type"`
		Title                    *string   `db:"title"`
		SettingsIsPrivate        bool      `db:"settings_is_private"`
		SettingsShuffleQuestions bool      `db:"settings_shuffle_questions"`
		SettingsShuffleAnswers   bool      `db:"settings_shuffle_answers"`
		CreatedAt                time.Time `db:"created_at"`
	}

	DefaultRepository struct {
		db *sqlx.DB
	}
)

func NewRepository(db *sqlx.DB) Repository {
	return &DefaultRepository{db: db}
}

func (r *DefaultRepository) Insert(ctx context.Context, tx transactional.Tx, in *model.Game) error {
	const query = `
		insert into game (id, status, "type", author_id, title) values ($1, $2, $3, $4, $5)
		on conflict (id) do nothing
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.Status, in.Type, in.AuthorID, in.Title)
	if err != nil {
		return err
	}

	const settingsQuery = `
		insert into game_settings (
	    	game_id,
			is_private,
		    shuffle_questions, 
		    shuffle_answers
		) values ($1, $2, $3, $4) 
		on conflict (game_id) do nothing
	`

	_, err = tx.ExecContext(
		ctx,
		settingsQuery,
		in.ID,
		in.Settings.IsPrivate,
		in.Settings.ShuffleQuestions,
		in.Settings.ShuffleAnswers,
	)
	return err
}

func (r *DefaultRepository) Update(ctx context.Context, tx transactional.Tx, in *model.Game) error {
	const query = `
		update game set 
			status = $2,
			updated_at = now()
		where id = $1
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.Status)
	return err
}

func (r *DefaultRepository) Get(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	return r.get(ctx, r.db, id)
}

func (r *DefaultRepository) GetWithTx(ctx context.Context, tx transactional.Tx, id uuid.UUID) (*model.Game, error) {
	return r.get(ctx, tx, id)
}

func (r *DefaultRepository) GetByAuthorID(ctx context.Context, authorID uuid.UUID) ([]model.Game, error) {
	const query = `
		select 
			g.id, 
			g."type", 
			g.status, 
			g.author_id,
			g.created_at,
			g.title,
			gs.is_private as settings_is_private, 
			gs.shuffle_questions as settings_shuffle_questions,
			gs.shuffle_answers as settings_shuffle_answers
		from game as g
		inner join game_settings as gs on gs.game_id = g.id
		where g.author_id = $1
	`

	var result []sqlxGame
	if err := r.db.SelectContext(ctx, &result, query, authorID); err != nil {
		return nil, err
	}

	return slices.SafeMap(result, func(i sqlxGame) model.Game {
		return convertToGame(i)
	}), nil
}

func (r *DefaultRepository) InsertGameQuestions(ctx context.Context, tx transactional.Tx, gameID uuid.UUID, questionIDs []uuid.UUID) error {
	const query = `
		insert into game_question (game_id, question_id) select $1, unnest($2::UUID[])
		on conflict (game_id, question_id) do nothing
	`

	_, err := tx.ExecContext(ctx, query, gameID, pq.Array(questionIDs))
	return err
}

func (r *DefaultRepository) GetQuestionIDsBySpec(ctx context.Context, tx transactional.Tx, spec *Spec) ([]uuid.UUID, error) {
	const query = `
		select question_id
		from game_question
		where game_id = $1
		and ($2::UUID[] is null or cardinality($2::UUID[]) = 0 or question_id != ANY($2::UUID[]))
	`

	var result []uuid.UUID
	if err := tx.SelectContext(ctx, &result, query, spec.GameID, pq.Array(spec.ExcludeQuestionIDs)); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *DefaultRepository) get(ctx context.Context, db transactional.Tx, id uuid.UUID) (*model.Game, error) {
	const query = `
		select 
			g.id, 
			g."type", 
			g.status, 
			g.author_id,
			g.created_at,
			g.title,
			gs.is_private as settings_is_private, 
			gs.shuffle_questions as settings_shuffle_questions,
			gs.shuffle_answers as settings_shuffle_answers
		from game as g
		inner join game_settings as gs on gs.game_id = g.id
		where g.id = $1
		limit 1
	`

	var result sqlxGame
	if err := db.GetContext(ctx, &result, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, contracts.ErrGameNotFound
		}

		return nil, err
	}

	return structs.Pointer(convertToGame(result)), nil
}

func convertToGame(in sqlxGame) model.Game {
	return model.Game{
		ID:       in.ID,
		Type:     model.GameType(in.Type),
		Status:   model.GameStatus(in.Status),
		AuthorID: in.AuthorID,
		Title:    in.Title,
		Settings: model.GameSettings{
			IsPrivate:        in.SettingsIsPrivate,
			ShuffleQuestions: in.SettingsShuffleQuestions,
			ShuffleAnswers:   in.SettingsShuffleAnswers,
		},
		CreatedAt: in.CreatedAt,
	}
}
