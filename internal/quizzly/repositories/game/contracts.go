package game

import (
	"context"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type (
	Spec struct {
		AuthorID  *uuid.UUID
		IsPrivate *bool
		Limit     int64
		Statuses  []model.GameStatus
	}

	QuestionSpec struct {
		GameID             uuid.UUID
		ExcludeQuestionIDs []uuid.UUID
	}

	Repository interface {
		Insert(ctx context.Context, tx transactional.Tx, in *model.Game) error
		Update(ctx context.Context, tx transactional.Tx, in *model.Game) error
		Get(ctx context.Context, id uuid.UUID) (*model.Game, error)
		GetWithTx(ctx context.Context, tx transactional.Tx, id uuid.UUID) (*model.Game, error)
		GetBySpec(ctx context.Context, spec *Spec) ([]model.Game, error)

		InsertGameQuestions(ctx context.Context, tx transactional.Tx, gameID uuid.UUID, questionIDs []uuid.UUID) error
		GetQuestionIDsBySpec(ctx context.Context, tx transactional.Tx, spec *QuestionSpec) ([]uuid.UUID, error)
	}
)
