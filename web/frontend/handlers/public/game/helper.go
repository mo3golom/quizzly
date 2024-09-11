package game

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func getResultsLink(gameID uuid.UUID, playerID uuid.UUID, request ...*http.Request) string {
	link := fmt.Sprintf("/game/results?id=%s&game_id=%s", playerID.String(), gameID.String())
	if len(request) > 0 {
		return fmt.Sprintf("%s%s", request[0].Host, link)
	}

	return link
}
