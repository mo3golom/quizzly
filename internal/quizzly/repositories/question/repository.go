package question

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	defaultLimit int64 = 10_000_000
)

type (
	sqlxQuestion struct {
		ID                    uuid.UUID            `db:"id"`
		ImageID               *string              `db:"image_id"`
		Text                  string               `db:"text"`
		Type                  string               `db:"type"`
		CreatedAt             time.Time            `db:"created_at"`
		AnswerOptionID        model.AnswerOptionID `db:"answer_option_id"`
		AnswerOptionAnswer    string               `db:"answer_option_answer"`
		AnswerOptionIsCorrect bool                 `db:"answer_option_is_correct"`
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

func (r *DefaultRepository) Insert(ctx context.Context, tx transactional.Tx, in *model.Question) error {
	const query = ` 
		insert into question (id, "text", "type", author_id, image_id) values ($1, $2, $3, $4, $5) on conflict (id) do nothing
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.Text, in.Type, in.AuthorID, in.ImageID)
	if err != nil {
		return err
	}

	const answerOptionsQuery = ` 
		insert into question_answer_option (question_id, answer, is_correct)
		select $1, unnest($2::text[]), unnest($3::boolean[])
	`

	answer := make([]string, 0, len(in.AnswerOptions))
	isCorrect := make([]bool, 0, len(in.AnswerOptions))
	for _, item := range in.AnswerOptions {
		answer = append(answer, item.Answer)
		isCorrect = append(isCorrect, item.IsCorrect)
	}

	_, err = tx.ExecContext(ctx, answerOptionsQuery, in.ID, pq.Array(answer), pq.Array(isCorrect))
	return err
}

func (r *DefaultRepository) Delete(ctx context.Context, tx transactional.Tx, id uuid.UUID) error {
	const query = `update question set deleted_at = now() where id = $1`

	_, err := tx.ExecContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func (r *DefaultRepository) GetBySpec(ctx context.Context, spec *Spec) (*GetBySpecOut, error) {
	total, err := r.getBySpecTotalCount(ctx, spec)
	if err != nil {
		return nil, err
	}

	out, err := r.getBySpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	return &GetBySpecOut{
		Result:     out,
		TotalCount: total,
	}, nil
}

func (r *DefaultRepository) getBySpec(ctx context.Context, spec *Spec) ([]model.Question, error) {
	query := buildBaseGetBySpecQuery("q.id, q.text, q.type, q.image_id, q.created_at, qao.id as answer_option_id, qao.answer as answer_option_answer, qao.is_correct as answer_option_is_correct")

	limit := defaultLimit
	offset := int64(0)
	if spec.Page != nil {
		limit = spec.Page.Limit
		offset = (spec.Page.Number - 1) * spec.Page.Limit
	}

	var result []sqlxQuestion
	if err := r.db.SelectContext(
		ctx,
		&result,
		query,
		pq.Array(spec.IDs),
		spec.AuthorID,
		len(spec.IDs) > 0,
		limit,
		offset,
	); err != nil {
		return nil, err
	}

	return convert(result), nil
}

func (r *DefaultRepository) getBySpecTotalCount(ctx context.Context, spec *Spec) (int64, error) {
	query := buildBaseGetBySpecQuery("count(distinct(q.id))")

	var result int64
	if err := r.db.GetContext(
		ctx,
		&result,
		query,
		pq.Array(spec.IDs),
		spec.AuthorID,
		len(spec.IDs) > 0,
		defaultLimit,
		0,
	); err != nil {
		return 0, err
	}

	return result, nil
}

func buildBaseGetBySpecQuery(fields string) string {
	return fmt.Sprintf(` 
        with question_ids as (
    		select id from question
    		where ($1::UUID[] is null or cardinality($1::UUID[]) = 0 or id = ANY($1::UUID[]))
			and ($2::UUID is null or author_id = $2::UUID)
			and ($3::bool = true or deleted_at is null)
			order by created_at desc
    		limit $4 
	    	offset $5
		)
		select %s
		from question as q
		inner join question_answer_option as qao on qao.question_id = q.id
        inner join question_ids on q.id = question_ids.id
	`, fields)
}

func convert(in []sqlxQuestion) []model.Question {
	out := make([]model.Question, 0, len(in))
	tempMap := make(map[uuid.UUID][]sqlxQuestion, len(in))
	for _, item := range in {
		if _, ok := tempMap[item.ID]; !ok {
			tempMap[item.ID] = make([]sqlxQuestion, 0, 2)
		}

		tempMap[item.ID] = append(tempMap[item.ID], item)
	}

	for _, items := range tempMap {
		var result *model.Question
		for _, item := range items {
			if result == nil {
				result = &model.Question{
					ID:            item.ID,
					Text:          item.Text,
					Type:          model.QuestionType(item.Type),
					ImageID:       item.ImageID,
					AnswerOptions: make([]model.AnswerOption, 0, len(in)),
					CreatedAt:     item.CreatedAt,
				}
			}

			result.AnswerOptions = append(result.AnswerOptions, model.AnswerOption{
				ID:        item.AnswerOptionID,
				Answer:    item.AnswerOptionAnswer,
				IsCorrect: item.AnswerOptionIsCorrect,
			})
		}

		if result == nil {
			continue
		}
		out = append(out, *result)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})

	return out
}
