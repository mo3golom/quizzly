package session

import (
	"context"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type (
	Spec struct {
		PlayerID uuid.UUID
		GameID   uuid.UUID
	}

	ItemSpec struct {
		PlayerID   uuid.UUID
		GameID     uuid.UUID
		QuestionID *uuid.UUID
	}

	Repository interface {
		Insert(ctx context.Context, tx transactional.Tx, in *model.Session) error
		Update(ctx context.Context, tx transactional.Tx, in *model.Session) error
		GetBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *Spec) (*model.Session, error)

		InsertSessionItem(ctx context.Context, tx transactional.Tx, in *model.SessionItem) error
		UpdateSessionItem(ctx context.Context, tx transactional.Tx, in *model.SessionItem) error
		DeleteSessionItemsBySessionID(ctx context.Context, tx transactional.Tx, sessionID int64) error
		GetSessionBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *ItemSpec) ([]model.SessionItem, error)
		GetSessionsByGameID(ctx context.Context, id uuid.UUID) ([]model.SessionExtended, error)
	}
)
