package supabase

import (
	"context"
	"github.com/google/uuid"
)

type DefaultAuthContext struct {
	context.Context
	userID uuid.UUID
}

func (d DefaultAuthContext) UserID() uuid.UUID {
	return d.userID
}
