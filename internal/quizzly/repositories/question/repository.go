package question

import (
	"context"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	sqlxQuestion struct {
		ID                    uuid.UUID `db:"id"`
		Text                  string    `db:"text"`
		Type                  string    `db:"type"`
		AnswerOptionAnswer    string    `db:"answer_option_answer"`
		AnswerOptionIsCorrect bool      `db:"answer_option_is_correct"`
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
		insert into question (id, "text", "type", author_id) values ($1, $2, $3, $4) on conflict (id) do nothing
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.Text, in.Type, in.AuthorID)
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

func (r *DefaultRepository) GetBySpec(ctx context.Context, spec *Spec) ([]model.Question, error) {
	const query = ` 
		select q.id, q.text, q.type, qao.answer as answer_option_answer, qao.is_correct as answer_option_is_correct 
		from question as q
		inner join question_answer_option as qao on qao.question_id = q.id
		where ($1::UUID[] is null or cardinality($1::UUID[]) = 0 or q.id = ANY($1::UUID[]))
		and ($2::UUID is null or q.author_id = $2::UUID)
	`

	var result []sqlxQuestion
	if err := r.db.SelectContext(
		ctx,
		&result,
		query,
		pq.Array(spec.IDs),
		spec.AuthorID,
	); err != nil {
		return nil, err
	}

	return convert(result), nil
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
					AnswerOptions: make([]model.AnswerOption, 0, len(in)),
				}
			}

			result.AnswerOptions = append(result.AnswerOptions, model.AnswerOption{
				Answer:    item.AnswerOptionAnswer,
				IsCorrect: item.AnswerOptionIsCorrect,
			})
		}

		if result == nil {
			continue
		}
		out = append(out, *result)
	}

	return out
}
