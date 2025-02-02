package game

import (
	"context"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"time"
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
)

func (r *DefaultRepository) Upsert(ctx context.Context, in *model.Game) error {
	const query = `
		insert into game (id, status, "type", author_id, title) values ($1, $2, $3, $4, $5)
		on conflict (id) do update set
			status = excluded.status,
			title = excluded.title
	`

	var title *string
	if in.Title != nil && len(*in.Title) > 0 {
		title = in.Title
	}

	_, err := r.db(ctx).ExecContext(ctx, query, in.ID, in.Status, in.Type, in.AuthorID, title)
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
		on conflict (game_id) do update set
			is_private = excluded.is_private,
			shuffle_questions = excluded.shuffle_questions,
			shuffle_answers = excluded.shuffle_answers,
			show_right_answers = excluded.show_right_answers,
			input_custom_name = excluded.input_custom_name
	`

	_, err = r.db(ctx).ExecContext(
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

func (r *DefaultRepository) GetBySpec(ctx context.Context, spec *Spec) ([]model.Game, error) {
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
		where ($1::UUID[] is null or cardinality($1::UUID[]) = 0 or g.id = any($1))
		  and ($2::UUID is null or g.author_id = $2)
		  and ($3::bool is null or gs.is_private = $3)
		  and ($4::text[] is null or cardinality($4::text[]) = 0 or g.status = any($4))
		order by g.created_at desc
		limit $5
	`

	limit := defaultLimit
	if spec.Limit > 0 {
		limit = spec.Limit
	}

	var result []sqlxGame
	if err := r.db(ctx).SelectContext(
		ctx,
		&result,
		query,
		pq.Array(spec.IDs),
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
