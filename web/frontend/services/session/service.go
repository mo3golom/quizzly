package session

import (
	"context"
	"github.com/a-h/templ"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	frontend_admin_game "quizzly/web/frontend/templ/admin/game"
)

type (
	DefaultService struct {
		sessions contracts.SessionUsecase
	}
)

func NewService(
	sessions contracts.SessionUsecase,
) Service {
	return &DefaultService{
		sessions: sessions,
	}
}

func (s *DefaultService) List(ctx context.Context, spec *Spec, _ *ListOptions) ([]templ.Component, error) {
	specificSessions, err := s.sessions.GetSessions(ctx, spec.GameID)
	if err != nil {
		return nil, err
	}

	return slices.SafeMap(specificSessions, func(session model.SessionExtended) templ.Component {
		return frontend_admin_game.SessionListItem(
			session.PlayerID.String(),
			int(session.CompletionRate()),
			session.Status,
		)
	}), nil
}
