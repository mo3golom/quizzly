package session

import (
	"context"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/handlers"
	frontend_admin_game "quizzly/web/frontend/templ/admin/game"
	"sort"
	"time"
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

func (s *DefaultService) List(ctx context.Context, spec *Spec, page int64, limit int64) (*ListOut, error) {
	specificSessions, err := s.sessions.GetExtendedSessions(ctx, spec.GameID, page, limit)
	if err != nil {
		return nil, err
	}

	specificPlayers, err := s.players.Get(
		ctx,
		slices.SafeMap(specificSessions.Result, func(session model.ExtendedSession) uuid.UUID {
			return session.PlayerID
		}),
	)
	if err != nil {
		return nil, err
	}

	specificPlayersMap := make(map[uuid.UUID]model.Player, len(specificPlayers))
	for _, player := range specificPlayers {
		specificPlayersMap[player.ID] = player
	}

	sort.Slice(specificSessions.Result, func(i, j int) bool {
		return specificSessions.Result[i].Session.ID > specificSessions.Result[j].ID
	})
	return &ListOut{
		Result: slices.SafeMap(specificSessions.Result, func(session model.ExtendedSession) templ.Component {
			sessionStartedAt := findSessionStart(session.Items)
			sessionLastQuestionAnsweredAt := findSessionLastAnswerTime(session.Items)
			moscowLocation, _ := time.LoadLocation("Europe/Moscow")
			if sessionStartedAt != nil {
				sessionStartedAt = structs.Pointer(sessionStartedAt.In(moscowLocation))
			}
			if sessionLastQuestionAnsweredAt != nil {
				sessionLastQuestionAnsweredAt = structs.Pointer(sessionLastQuestionAnsweredAt.In(moscowLocation))
			}

			return frontend_admin_game.SessionListItem(handlers.SessionItemStatistics{
				PlayerName:                    specificPlayersMap[session.PlayerID].Name,
				CompletionRate:                int(session.CompletionRate()),
				SessionStatus:                 session.Status,
				SessionStartedAt:              sessionStartedAt,
				SessionLastQuestionAnsweredAt: sessionLastQuestionAnsweredAt,
			},
			)
		}),
		TotalCount: specificSessions.TotalCount,
	}, nil
}

func findSessionStart(in []model.SessionItem) *time.Time {
	var minimum *time.Time
	for _, item := range in {
		if minimum == nil {
			minimum = &item.CreatedAt
			continue
		}

		if item.CreatedAt.Before(*minimum) {
			minimum = &item.CreatedAt
		}
	}

	return minimum
}

func findSessionLastAnswerTime(in []model.SessionItem) *time.Time {
	var maximum *time.Time
	for _, item := range in {
		if maximum == nil {
			maximum = item.AnsweredAt
			continue
		}

		if item.AnsweredAt != nil && item.AnsweredAt.After(*maximum) {
			maximum = &item.CreatedAt
		}
	}

	return maximum
}
