package question

import (
	"context"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type Usecase struct {
	questions question.Repository
	template  transactional.Template
}

func NewUsecase(
	questions question.Repository,
	template transactional.Template,
) contracts.QuestionUsecase {
	return &Usecase{
		questions: questions,
		template:  template,
	}
}

func (u *Usecase) Create(ctx context.Context, in *model.Question) error {
	if len(in.AnswerOptions) == 0 {
		return contracts.ErrEmptyAnswerOptions
	}

	if in.ID == uuid.Nil {
		in.ID = uuid.New()
	}

	if in.Type == model.QuestionTypeFillTheGap {
		in.AnswerOptions[0].IsCorrect = true
	}

	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.questions.Insert(ctx, tx, in)
	})
}

func (u *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.questions.Delete(ctx, tx, id)
	})
}

func (u *Usecase) GetByAuthor(ctx context.Context, authorID uuid.UUID, page int64, limit int64) (*contracts.GetByAuthorOut, error) {
	result, err := u.questions.GetBySpec(ctx, &question.Spec{
		AuthorID: &authorID,
		Page: &question.Page{
			Number: page,
			Limit:  limit,
		},
	})
	if err != nil {
		return nil, err
	}

	return &contracts.GetByAuthorOut{
		Result:     result.Result,
		TotalCount: result.TotalCount,
	}, nil
}

func (u *Usecase) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Question, error) {
	result, err := u.questions.GetBySpec(ctx, &question.Spec{
		IDs: ids,
	})
	if err != nil {
		return nil, err
	}

	return result.Result, nil
}
