package game

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"quizzly/internal/quizzly/model"
	"time"
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
)

func (r *DefaultRepository) InsertQuestion(ctx context.Context, in *model.Question) error {
	const query = ` 
        with last_question_sort as (
		    select game_id, sort
		    from question
		    where game_id = $4
		    order by sort desc
		    limit 1
		),
		data as (
		    select
		        $1::uuid as id,
		        $2 as text, 
		        $3 as type ,
		        $4::uuid as game_id,
		        $5 as image_id
		)
		insert into question (id, "text", "type", "game_id", "image_id", "sort")
		select d.id, d.text, d.type, d.game_id, d.image_id, coalesce(lqs.sort, 0) + 1 as sort
		from data as d
		left join last_question_sort lqs on lqs.game_id = d.game_id
	`

	_, err := r.db(ctx).ExecContext(ctx, query, in.ID, in.Text, in.Type, in.GameID, in.ImageID)
	if err != nil {
		return err
	}

	return r.upsertAnswerOption(ctx, in)
}

func (r *DefaultRepository) UpdateQuestion(ctx context.Context, in *model.Question) error {
	const query = `
		update question set 
			"text" = $2,
			image_id = $3     
		where id = $1
	`
	_, err := r.db(ctx).ExecContext(ctx, query, in.ID, in.Text, in.ImageID)
	if err != nil {
		return err
	}

	return r.upsertAnswerOption(ctx, in)
}

func (r *DefaultRepository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	const query = `update question set deleted_at = now() where id = $1`

	_, err := r.db(ctx).ExecContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func (r *DefaultRepository) upsertAnswerOption(ctx context.Context, in *model.Question) error {
	const answerOptionsQueryDelete = `delete from question_answer_option where question_id = $1`
	_, err := r.db(ctx).ExecContext(ctx, answerOptionsQueryDelete, in.ID)
	if err != nil {
		return err
	}

	const answerOptionsQueryInsert = ` 
		insert into question_answer_option (question_id, answer, is_correct)
		select $1, unnest($2::text[]), unnest($3::boolean[])
	`

	answer := make([]string, 0, len(in.AnswerOptions))
	isCorrect := make([]bool, 0, len(in.AnswerOptions))
	for _, item := range in.AnswerOptions {
		answer = append(answer, item.Answer)
		isCorrect = append(isCorrect, item.IsCorrect)
	}

	_, err = r.db(ctx).ExecContext(ctx, answerOptionsQueryInsert, in.ID, pq.Array(answer), pq.Array(isCorrect))
	return err
}

func (r *DefaultRepository) GetQuestionsBySpec(ctx context.Context, spec *QuestionsSpec) ([]model.Question, error) {
	const query = `
       select 
           q.id, 
           q.text, 
           q.type, 
           q.image_id, 
           q.created_at,
           qao.id as answer_option_id, 
           qao.answer as answer_option_answer, 
           qao.is_correct as answer_option_is_correct
		from question as q
		inner join question_answer_option as qao on qao.question_id = q.id
        where ($1::UUID[] is null or cardinality($1::UUID[]) = 0 or q.id = ANY($1::UUID[]))
		  and ($2::UUID is null or game_id = $2::UUID)
		  and deleted_at is null
       order by q.sort, q.created_at
	`

	var result []sqlxQuestion
	if err := r.db(ctx).SelectContext(
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
			ID:        item.AnswerOptionID,
			Answer:    item.AnswerOptionAnswer,
			IsCorrect: item.AnswerOptionIsCorrect,
		})
	}

	return out
}
