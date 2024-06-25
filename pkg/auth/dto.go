package auth

import (
	"github.com/google/uuid"
	"time"
)

type (
	insertLoginCodeIn struct {
		code      LoginCode
		userID    uuid.UUID
		expiresAt time.Time
	}

	getLoginCodeIn struct {
		code   LoginCode
		userID uuid.UUID
	}

	loginCodeExtended struct {
		userID uuid.UUID
		code   LoginCode
	}

	user struct {
		id    uuid.UUID
		email Email
	}
)
