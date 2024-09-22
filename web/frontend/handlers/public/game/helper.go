package game

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func gameLink(gameID uuid.UUID, request ...*http.Request) string {
	link := fmt.Sprintf("/game/%s", gameID.String())
	if len(request) > 0 {
		return fmt.Sprintf("%s%s", request[0].Host, link)
	}

	return link
}

func resultsLink(gameID uuid.UUID, playerID uuid.UUID, request ...*http.Request) string {
	link := fmt.Sprintf("/game/%s/results/%s", gameID.String(), playerID.String())
	if len(request) > 0 {
		return fmt.Sprintf("%s%s", request[0].Host, link)
	}

	return link
}
