package link

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"quizzly/pkg/variables"
)

type DefaultService struct {
	variables variables.Repository
}

func NewService(variables variables.Repository) Service {
	return &DefaultService{
		variables: variables,
	}
}

func (s *DefaultService) GameLink(gameID uuid.UUID, request ...*http.Request) string {
	link := fmt.Sprintf("/game/%s", gameID.String())

	return newLinkBuilder(link).
		addHost(request...).
		addHTTPS(s.variables).
		build()
}

func (s *DefaultService) GameResultsLink(gameID uuid.UUID, playerID uuid.UUID, request ...*http.Request) string {
	link := fmt.Sprintf("/game/%s/results/%s", gameID.String(), playerID.String())

	return newLinkBuilder(link).
		addHost(request...).
		addHTTPS(s.variables).
		build()
}
