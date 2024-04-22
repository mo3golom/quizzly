package contracts

import (
	"context"
	"quizzly/internal/quizzly/model"
)

type (
	QuestionUsecase interface {
		Create(ctx context.Context, in *model.Question) error
	}
)
