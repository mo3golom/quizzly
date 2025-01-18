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

	GetExtendedSessionSpec struct {
		GameID uuid.UUID
		Page   *Page
	}

	GetExtendedSessionsBySpecOut struct {
		Result     []model.ExtendedSession
		TotalCount int64
	}

	Page struct {
		Number int64
		Limit  int64
	}

	Repository interface {
		Insert(ctx context.Context, tx transactional.Tx, in *model.Session) error
		Update(ctx context.Context, tx transactional.Tx, in *model.Session) error
		GetBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *Spec) (*model.Session, error)

		InsertSessionItem(ctx context.Context, tx transactional.Tx, in *model.SessionItem) error
		DeleteSessionItemsBySessionID(ctx context.Context, tx transactional.Tx, sessionID int64) error
		GetSessionBySpecWithTx(ctx context.Context, tx transactional.Tx, spec *ItemSpec) ([]model.SessionItem, error)
		GetExtendedSessionsBySpec(ctx context.Context, spec *GetExtendedSessionSpec) (*GetExtendedSessionsBySpecOut, error)
	}
)
