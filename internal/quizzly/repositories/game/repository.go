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

const (
	defaultLimit int64 = 10_000_000
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
		SettingsShowRightAnswers bool      `db:"settings_show_right_answers"`
		SettingsInputCustomName  bool      `db:"settings_input_custom_name"`
		CreatedAt                time.Time `db:"created_at"`
	}

	sqlxGameQuestion struct {
		ID   uuid.UUID `db:"question_id"`
		Sort *int64    `db:"sort"`
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

	var title *string
	if in.Title != nil && len(*in.Title) > 0 {
		title = in.Title
	}

	_, err := tx.ExecContext(ctx, query, in.ID, in.Status, in.Type, in.AuthorID, title)
	if err != nil {
		return err
	}

	const settingsQuery = `
		insert into game_settings (
	    	game_id,
			is_private,
		    shuffle_questions, 
		    shuffle_answers,
		    show_right_answers,
		    input_custom_name
		) values ($1, $2, $3, $4, $5, $6) 
		on conflict (game_id) do nothing
	`

	_, err = tx.ExecContext(
		ctx,
		settingsQuery,
		in.ID,
		in.Settings.IsPrivate,
		in.Settings.ShuffleQuestions,
		in.Settings.ShuffleAnswers,
		in.Settings.ShowRightAnswers,
		in.Settings.InputCustomName,
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

func (r *DefaultRepository) GetBySpec(ctx context.Context, spec *Spec) ([]model.Game, error) {
	if spec == nil || (spec.AuthorID == nil && spec.IsPrivate == nil) {
		return nil, nil
	}

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
			gs.shuffle_answers as settings_shuffle_answers,
			gs.show_right_answers as settings_show_right_answers,
		    gs.input_custom_name as settings_input_custom_name
		from game as g
		inner join game_settings as gs on gs.game_id = g.id
		where ($1::UUID is null or g.author_id = $1)
		  and ($2::bool is null or gs.is_private = $2)
		  and ($3::text[] is null or cardinality($3::text[]) = 0 or g.status = any($3))
		order by g.created_at desc
		limit $4
	`

	limit := defaultLimit
	if spec.Limit > 0 {
		limit = spec.Limit
	}

	var result []sqlxGame
	if err := r.db.SelectContext(
		ctx,
		&result,
		query,
		spec.AuthorID,
		spec.IsPrivate,
		pq.Array(spec.Statuses),
		limit,
	); err != nil {
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

func (r *DefaultRepository) GetQuestionIDsBySpec(ctx context.Context, tx transactional.Tx, spec *QuestionSpec) ([]GameQuestion, error) {
	const query = `
		select question_id, sort
		from game_question
		where game_id = $1
		and ($2::UUID[] is null or cardinality($2::UUID[]) = 0 or question_id != ANY($2::UUID[]))
	`

	var result []sqlxGameQuestion
	if err := tx.SelectContext(ctx, &result, query, spec.GameID, pq.Array(spec.ExcludeQuestionIDs)); err != nil {
		return nil, err
	}

	return slices.SafeMap(result, func(i sqlxGameQuestion) GameQuestion {
		var sort int64
		if i.Sort != nil {
			sort = *i.Sort
		}

		return GameQuestion{
			ID:   i.ID,
			Sort: sort,
		}
	}), nil
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
			gs.shuffle_answers as settings_shuffle_answers,
			gs.show_right_answers as settings_show_right_answers,
		    gs.input_custom_name as settings_input_custom_name
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
			ShowRightAnswers: in.SettingsShowRightAnswers,
			InputCustomName:  in.SettingsInputCustomName,
		},
		CreatedAt: in.CreatedAt,
	}
}
