package game

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func getGameLink(gameID uuid.UUID, request *http.Request) string {
	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/game/play?id=%s", scheme, request.Host, gameID.String())
}
