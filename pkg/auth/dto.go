package auth

import (
	"time"

	"github.com/google/uuid"
)

type (
	upsertLoginCodeIn struct {
		code      LoginCode
		userID    uuid.UUID
		expiresAt time.Time
	}

	upsertTokenIn struct {
		token     Token
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
