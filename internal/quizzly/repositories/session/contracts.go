package session

import (
	"context"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/model"
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
		Insert(ctx context.Context, in *model.Session) error
		Update(ctx context.Context, in *model.Session) error
		GetBySpec(ctx context.Context, spec *Spec) (*model.Session, error)

		InsertSessionItem(ctx context.Context, in *model.SessionItem) error
		DeleteSessionItemsBySessionID(ctx context.Context, sessionID int64) error
		GetSessionBySpec(ctx context.Context, spec *ItemSpec) ([]model.SessionItem, error)
		GetExtendedSessionsBySpec(ctx context.Context, spec *GetExtendedSessionSpec) (*GetExtendedSessionsBySpecOut, error)
	}
)
