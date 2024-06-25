package session

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"
	"time"

	"github.com/google/uuid"
)

type (
	sqlxSession struct {
		ID       int64     `db:"id"`
		PlayerID uuid.UUID `db:"player_id"`
		GameID   uuid.UUID `db:"game_id"`
		Status   string    `db:"status"`
	}

	sqlxSessionExtended struct {
		ID         int64      `db:"id"`
		PlayerID   uuid.UUID  `db:"player_id"`
		GameID     uuid.UUID  `db:"game_id"`
		Status     string     `db:"status"`
		ItemID     *int64     `db:"item_id"`
		QuestionID *uuid.UUID `db:"item_question_id"`
		Answers    []byte     `db:"item_answers"`
		IsCorrect  *bool      `db:"item_is_correct"`
		AnsweredAt *time.Time `db:"item_answered_at"`
	}

	sqlxSessionItem struct {
		ID         int64      `db:"id"`
		SessionID  int64      `db:"session_id"`
		QuestionID uuid.UUID  `db:"question_id"`
		Answers    []byte     `db:"answers"`
		IsCorrect  *bool      `db:"is_correct"`
		AnsweredAt *time.Time `db:"answered_at"`
	}

	DefaultRepository struct {
		db *sqlx.DB
	}
)

func NewRepository(db *sqlx.DB) Repository {
	return &DefaultRepository{db: db}
}

func (r *DefaultRepository) Insert(ctx context.Context, tx transactional.Tx, in *model.Session) error {
	const query = `
		insert into player_session (game_id, player_id, status) values ($1, $2, $3)
	`

	_, err := tx.ExecContext(ctx, query, in.GameID, in.PlayerID, in.Status)
	return err
}

func (r *DefaultRepository) Update(ctx context.Context, tx transactional.Tx, in *model.Session) error {
	const query = `
		update player_session set 
			status = $2,
			updated_at = now()
		where id = $1
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.Status)
	return err
}

func (r *DefaultRepository) GetBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *Spec) (*model.Session, error) {
	const query = `
		select id, game_id, player_id, status 
		from player_session 
		where player_id = $1 and game_id = $2
		limit 1
	`

	var result sqlxSession
	if err := tx.GetContext(ctx, &result, query, spec.PlayerID, spec.GameID); err != nil {
		return nil, err
	}

	return &model.Session{
		ID:       result.ID,
		GameID:   result.GameID,
		PlayerID: result.PlayerID,
		Status:   model.SessionStatus(result.Status),
	}, nil
}

func (r *DefaultRepository) InsertSessionItem(ctx context.Context, tx transactional.Tx, in *model.SessionItem) error {
	const query = `
		insert into player_session_item (session_id, question_id) values ($1, $2)
	`

	_, err := tx.ExecContext(ctx, query, in.SessionID, in.QuestionID)
	return err
}

func (r *DefaultRepository) UpdateSessionItem(ctx context.Context, tx transactional.Tx, in *model.SessionItem) error {
	const query = `
		update player_session_item set 
			answers = $2,
			is_correct = $3,
			answered_at = $4,
			updated_at = now()
		where id = $1
	`

	answers, err := json.Marshal(in.Answers)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, in.ID, answers, in.IsCorrect, in.AnsweredAt)
	return err
}

func (r *DefaultRepository) GetSessionBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *ItemSpec) ([]model.SessionItem, error) {
	const query = `
		select psi.id, psi.session_id, psi.question_id, psi.answers, psi.is_correct, psi.answered_at
		from player_session_item as psi
		inner join player_session as ps on ps.id = psi.session_id
		where ps.player_id = $1 
	      and ps.game_id = $2
	      and ($3::UUID is null or psi.question_id = $3::UUID)
	`

	var result []sqlxSessionItem
	if err := tx.SelectContext(ctx, &result, query, spec.PlayerID, spec.GameID, spec.QuestionID); err != nil {
		return nil, err
	}

	return slices.Map(result, func(i sqlxSessionItem) (model.SessionItem, error) {
		out, err := convertSessionItem(&i)
		if err != nil {
			return model.SessionItem{}, err
		}
		return *out, nil
	})
}

func (r *DefaultRepository) GetSessionsByGameID(ctx context.Context, id uuid.UUID) ([]model.SessionExtended, error) {
	const query = `
		select 
		    ps.id, 
		    ps.game_id, 
		    ps.player_id, 
		    ps.status, 
		    psi.id as item_id, 
		    psi.question_id as item_question_id, 
		    psi.answers as item_answers, 
		    psi.is_correct as item_is_correct, 
		    psi.answered_at as item_answered_at
		from player_session ps
		left join player_session_item as psi on psi.session_id = ps.id
		where ps.game_id = $1
	`

	var result []sqlxSessionExtended
	if err := r.db.SelectContext(ctx, &result, query, id); err != nil {
		return nil, err
	}

	resultMap := make(map[int64]model.SessionExtended, len(result))
	for _, item := range result {
		session, ok := resultMap[item.ID]
		if !ok {
			session = model.SessionExtended{
				Session: model.Session{
					ID:       item.ID,
					GameID:   item.GameID,
					PlayerID: item.PlayerID,
					Status:   model.SessionStatus(item.Status),
				},
				Items: make([]model.SessionItem, 0, 1),
			}
		}

		if item.ItemID != nil {
			sessionItems := resultMap[item.ID].Items
			sessionItems = append(sessionItems, model.SessionItem{
				ID:         *item.ItemID,
				QuestionID: *item.QuestionID,
				IsCorrect:  item.IsCorrect,
				AnsweredAt: item.AnsweredAt,
			})
			session.Items = sessionItems
		}

		resultMap[item.ID] = session
	}

	out := make([]model.SessionExtended, 0, len(resultMap))
	for _, item := range resultMap {
		item := item
		out = append(out, item)
	}

	return out, nil
}

func convertSessionItem(in *sqlxSessionItem) (*model.SessionItem, error) {
	var answers []string
	if in.Answers != nil {
		if err := json.Unmarshal(in.Answers, &answers); err != nil {
			return nil, err
		}
	}

	return &model.SessionItem{
		ID:         in.ID,
		SessionID:  in.SessionID,
		QuestionID: in.QuestionID,
		Answers:    answers,
		IsCorrect:  in.IsCorrect,
		AnsweredAt: in.AnsweredAt,
	}, nil
}
