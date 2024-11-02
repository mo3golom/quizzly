package link

import (
	"github.com/google/uuid"
	"net/http"
)

type (
	Service interface {
		GameLink(gameID uuid.UUID, request ...*http.Request) string
		GameResultsLink(gameID uuid.UUID, playerID uuid.UUID, request ...*http.Request) string
	}
)
