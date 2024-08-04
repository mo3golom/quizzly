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

func (u *Usecase) GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]model.Question, error) {
	return u.questions.GetBySpec(ctx, &question.Spec{
		AuthorID: &authorID,
	})
}

func (u *Usecase) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Question, error) {
	return u.questions.GetBySpec(ctx, &question.Spec{
		IDs: ids,
	})
}
