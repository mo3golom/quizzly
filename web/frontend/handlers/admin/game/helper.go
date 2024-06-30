package game

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func getGameLink(gameID uuid.UUID, request *http.Request) string {
	return fmt.Sprintf("%s/game/play?id=%s", request.Host, gameID.String())
}
