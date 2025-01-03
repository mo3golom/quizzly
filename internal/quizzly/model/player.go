package model

import "github.com/google/uuid"

type (
	Player struct {
		ID              uuid.UUID
		Name            string
		NameUserEntered bool
		UserID          *uuid.UUID
	}
)
