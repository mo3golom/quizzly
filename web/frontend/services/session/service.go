package session

import (
	"fmt"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/auth"
	"quizzly/pkg/structs"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/handlers"
	frontend_admin_game "quizzly/web/frontend/templ/admin/game"
	"sort"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
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

func (s *DefaultService) List(request *http.Request, spec *Spec, page int64, limit int64) (*ListOut, error) {
	var currentPlayer *model.Player
	if authCtx, ok := request.Context().(auth.Context); ok && authCtx.UserID() != uuid.Nil {
		players, err := s.players.GetByUsers(request.Context(), []uuid.UUID{authCtx.UserID()})
		if err != nil {
			return nil, err
		}

		if len(players) > 0 {
			currentPlayer = structs.Pointer(players[0])
		}
	}

	specificSessions, err := s.sessions.GetExtendedSessions(request.Context(), spec.GameID, page, limit)
	if err != nil {
		return nil, err
	}

	specificPlayers, err := s.players.Get(
		request.Context(),
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
			moscowLocation, _ := time.LoadLocation("Europe/Moscow")
			sessionStartedAt := session.CreatedAt.In(moscowLocation)
			sessionLastQuestionAnsweredAt := findSessionLastAnswerTime(session.Items)
			if sessionLastQuestionAnsweredAt != nil {
				sessionLastQuestionAnsweredAt = structs.Pointer(sessionLastQuestionAnsweredAt.In(moscowLocation))
			}

			playerName := specificPlayersMap[session.PlayerID].Name
			if currentPlayer != nil && currentPlayer.Name == playerName {
				playerName = fmt.Sprintf("%s ( вы )", playerName)
			}

			return frontend_admin_game.SessionListItem(handlers.SessionItemStatistics{
				PlayerName:                    playerName,
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
