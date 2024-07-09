package session

import (
	"context"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	frontend_admin_game "quizzly/web/frontend/templ/admin/game"
	"sort"
)

type (
	DefaultService struct {
		sessions contracts.SessionUsecase
		players  contracts.PLayerUsecase
	}
)

func NewService(
	sessions contracts.SessionUsecase,
	players contracts.PLayerUsecase,
) Service {
	return &DefaultService{
		sessions: sessions,
		players:  players,
	}
}

func (s *DefaultService) List(ctx context.Context, spec *Spec, _ *ListOptions) ([]templ.Component, error) {
	specificSessions, err := s.sessions.GetSessions(ctx, spec.GameID)
	if err != nil {
		return nil, err
	}

	specificPlayers, err := s.players.Get(
		ctx,
		slices.SafeMap(specificSessions, func(session model.SessionExtended) uuid.UUID {
			return session.PlayerID
		}),
	)
	if err != nil {
		return nil, err
	}

	specificPlayersMap := make(map[uuid.UUID]model.Player, len(specificPlayers))
	for _, player := range specificPlayers {
		player = player
		specificPlayersMap[player.ID] = player
	}

	sort.Slice(specificSessions, func(i, j int) bool {
		return specificSessions[i].ID < specificSessions[j].ID
	})
	return slices.SafeMap(specificSessions, func(session model.SessionExtended) templ.Component {
		return frontend_admin_game.SessionListItem(
			specificPlayersMap[session.PlayerID].Name,
			int(session.CompletionRate()),
			session.Status,
		)
	}), nil
}
