package game

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"
	"time"
)

type (
	sqlxQuestion struct {
		ID                         uuid.UUID            `db:"id"`
		ImageID                    *string              `db:"image_id"`
		Text                       string               `db:"text"`
		Type                       string               `db:"type"`
		CreatedAt                  time.Time            `db:"created_at"`
		AnswerOptionID             model.AnswerOptionID `db:"answer_option_id"`
		AnswerOptionAnswer         string               `db:"answer_option_answer"`
		AnswerOptionIsCorrect      bool                 `db:"answer_option_is_correct"`
		AnswerOptionNextQuestionID *uuid.UUID           `db:"answer_option_next_question_id"`
	}
)

func (r *DefaultRepository) InsertQuestion(ctx context.Context, tx transactional.Tx, in *model.Question) error {
	const query = ` 
		insert into question (id, "text", "type", game_id, image_id) values ($1, $2, $3, $4, $5) on conflict (id) do nothing
	`

	_, err := tx.ExecContext(ctx, query, in.ID, in.Text, in.Type, in.GameID, in.ImageID)
	if err != nil {
		return err
	}

	return r.upsertAnswerOption(ctx, tx, in)
}

func (r *DefaultRepository) UpdateQuestion(ctx context.Context, tx transactional.Tx, in *model.Question) error {
	const query = `
		update question set 
			"text" = $2,
			image_id = $3
		where id = $1
	`
	_, err := tx.ExecContext(ctx, query, in.ID, in.Text, in.ImageID)
	if err != nil {
		return err
	}

	return r.upsertAnswerOption(ctx, tx, in)
}

func (r *DefaultRepository) DeleteQuestion(ctx context.Context, tx transactional.Tx, id uuid.UUID) error {
	const query = `update question set deleted_at = now() where id = $1`

	_, err := tx.ExecContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func (r *DefaultRepository) upsertAnswerOption(ctx context.Context, tx transactional.Tx, in *model.Question) error {
	const answerOptionsQueryDelete = `delete from question_answer_option where question_id = $1`
	_, err := tx.ExecContext(ctx, answerOptionsQueryDelete, in.ID)
	if err != nil {
		return err
	}

	const answerOptionsQueryInsert = ` 
		insert into question_answer_option (question_id, answer, is_correct, next_question_id)
		select $1, unnest($2::text[]), unnest($3::boolean[]), unnest($4::uuid[])
	`

	answer := make([]string, 0, len(in.AnswerOptions))
	isCorrect := make([]bool, 0, len(in.AnswerOptions))
	nextQuestionIDs := make([]*uuid.UUID, 0, len(in.AnswerOptions))
	for _, item := range in.AnswerOptions {
		answer = append(answer, item.Answer)
		isCorrect = append(isCorrect, item.IsCorrect)
		nextQuestionIDs = append(nextQuestionIDs, item.NextQuestionID)
	}

	_, err = tx.ExecContext(ctx, answerOptionsQueryInsert, in.ID, pq.Array(answer), pq.Array(isCorrect), pq.Array(nextQuestionIDs))
	return err
}

func (r *DefaultRepository) GetQuestionsBySpec(ctx context.Context, spec *QuestionsSpec) ([]model.Question, error) {
	return r.getQuestionsBySpec(ctx, r.db, spec)
}
func (r *DefaultRepository) GetQuestionsBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *QuestionsSpec) ([]model.Question, error) {
	return r.getQuestionsBySpec(ctx, tx, spec)
}

func (r *DefaultRepository) getQuestionsBySpec(ctx context.Context, tx transactional.Tx, spec *QuestionsSpec) ([]model.Question, error) {
	order := "q.created_at desc"
	if spec.Order != nil {
		order = fmt.Sprintf("q.%s %s", spec.Order.Field, spec.Order.Direction)
	}

	query := fmt.Sprintf(`
       select 
           q.id, 
           q.text, 
           q.type, 
           q.image_id, 
           q.created_at, 
           qao.id as answer_option_id, 
           qao.answer as answer_option_answer, 
           qao.is_correct as answer_option_is_correct,
           qao.next_question_id as answer_option_next_question_id
		from question as q
		inner join question_answer_option as qao on qao.question_id = q.id
        where ($1::UUID[] is null or cardinality($1::UUID[]) = 0 or q.id = ANY($1::UUID[]))
		  and ($2::UUID is null or game_id = $2::UUID)
		  and deleted_at is null
       order by %s
	`, order)

	var result []sqlxQuestion
	if err := tx.SelectContext(
		ctx,
		&result,
		query,
		pq.Array(spec.IDs),
		spec.GameID,
	); err != nil {
		return nil, err
	}

	return convertQuestions(result), nil
}

func convertQuestions(in []sqlxQuestion) []model.Question {
	out := make([]model.Question, 0, len(in))
	indexMap := make(map[uuid.UUID]int, len(in))
	for _, item := range in {
		index, ok := indexMap[item.ID]
		if !ok {
			out = append(out, model.Question{
				ID:            item.ID,
				Text:          item.Text,
				Type:          model.QuestionType(item.Type),
				ImageID:       item.ImageID,
				AnswerOptions: make([]model.AnswerOption, 0, 4),
				CreatedAt:     item.CreatedAt,
			})
			index = len(out) - 1
			indexMap[item.ID] = index
		}

		out[index].AnswerOptions = append(out[index].AnswerOptions, model.AnswerOption{
			ID:             item.AnswerOptionID,
			Answer:         item.AnswerOptionAnswer,
			IsCorrect:      item.AnswerOptionIsCorrect,
			NextQuestionID: item.AnswerOptionNextQuestionID,
		})
	}

	return out
}
