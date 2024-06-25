package game

import (
	"context"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type (
	sqlxGame struct {
		ID                       uuid.UUID `db:"id"`
		Status                   string    `db:"status"`
		Type                     string    `db:"type"`
		SettingsIsPrivate        bool      `db:"settings_is_private"`
		SettingsShuffleQuestions bool      `db:"settings_shuffle_questions"`
		SettingsShuffleAnswers   bool      `db:"settings_shuffle_answers"`
	}

	DefaultRepository struct {
	}
)

func NewRepository() Repository {
	return &DefaultRepository{}
}

func (r *DefaultRepository) Insert(ctx context.Context, tx transactional.Tx, in *model.Game) error {
	const query = `
		insert into game (id, status, "type") values ($1, $2, $3)
		on conflict (id) do nothing
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.Status, in.Type)
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

func (r *DefaultRepository) GetWithTx(ctx context.Context, tx transactional.Tx, id uuid.UUID) (*model.Game, error) {
	const query = `
		select 
			g.id, 
			g."type", 
			g.status, 
			gs.is_private as settings_is_private, 
			gs.shuffle_questions as settings_shuffle_questions,
			gs.shuffle_answers as settings_shuffle_answers
		from game as g
		inner join game_settings as gs on gs.game_id = g.id
		where g.id = $1
		limit 1
	`

	var result sqlxGame
	if err := tx.GetContext(ctx, &result, query, id); err != nil {
		return nil, err
	}

	return &model.Game{
		ID:     result.ID,
		Type:   model.GameType(result.Type),
		Status: model.GameStatus(result.Status),
		Settings: model.GameSettings{
			IsPrivate:        result.SettingsIsPrivate,
			ShuffleQuestions: result.SettingsShuffleQuestions,
			ShuffleAnswers:   result.SettingsShuffleAnswers,
		},
	}, nil
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
