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
	if in.ID == uuid.Nil {
		in.ID = uuid.New()
	}

	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.questions.Insert(ctx, tx, in)
	})
}
