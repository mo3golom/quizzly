package player

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/cookie"
	"quizzly/pkg/logger"
	"quizzly/pkg/structs"
	"quizzly/pkg/supabase"
	"time"
)

const (
	cookiePlayerID = "player-game"
)

type DefaultService struct {
	playerUC contracts.PLayerUsecase
	cookie   cookie.Service
	log      logger.Logger
}

func NewService(playerUC contracts.PLayerUsecase, cookieService cookie.Service, log logger.Logger) *DefaultService {
	return &DefaultService{
		playerUC: playerUC,
		cookie:   cookieService,
		log:      log,
	}
}

func (s *DefaultService) GetPlayer(writer http.ResponseWriter, request *http.Request, gameID uuid.UUID, customName ...string) (*model.Player, error) {
	var userID *uuid.UUID
	if authCtx, ok := request.Context().(supabase.AuthContext); ok && authCtx.UserID() != uuid.Nil {
		userID = structs.Pointer(authCtx.UserID())
	}

	name := ""
	if len(customName) > 0 {
		name = customName[0]
	}

	player, err := s.findPLayerID(request, userID, gameID, name)
	if err != nil {
		s.log.Error("error getting player id", err)
	}
	if player != nil {
		return player, s.setPlayerID(writer, player.ID, gameID)
	}

	player, err = s.newPlayer(request.Context(), userID, name)
	if err != nil {
		return nil, err
	}

	return player, s.setPlayerID(writer, player.ID, gameID)
}

func (s *DefaultService) findPLayerID(request *http.Request, userID *uuid.UUID, gameID uuid.UUID, customName string) (*model.Player, error) {
	player, err := s.findByUserID(request.Context(), userID)
	if err != nil {
		return nil, err
	}
	if player != nil {
		needUpdate := false
		if customName != "" {
			player.Name = customName
			player.NameUserEntered = true
			needUpdate = true
		}

		if needUpdate {
			err = s.playerUC.Update(request.Context(), player)
			if err != nil {
				return player, err
			}
		}

		return player, nil
	}

	player, err = s.findFromCookie(request, gameID)
	if err != nil {
		return nil, err
	}
	if player != nil {
		needUpdate := false
		if player.UserID == nil && userID != nil {
			player.UserID = userID
			needUpdate = true
		}
		if customName != "" {
			player.Name = customName
			player.NameUserEntered = true
			needUpdate = true
		}

		if needUpdate {
			err = s.playerUC.Update(request.Context(), player)
			if err != nil {
				return player, err
			}
		}

		return player, nil
	}

	return nil, nil
}

func (s *DefaultService) newPlayer(ctx context.Context, userID *uuid.UUID, customName string) (*model.Player, error) {
	name := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
	nameUserEntered := false
	if customName != "" {
		name = customName
		nameUserEntered = true
	}

	newPlayer := &model.Player{
		ID:              uuid.New(),
		Name:            name,
		NameUserEntered: nameUserEntered,
		UserID:          userID,
	}
	err := s.playerUC.Create(ctx, newPlayer)
	if err != nil {
		return nil, err
	}

	return newPlayer, nil
}

func (s *DefaultService) findByUserID(ctx context.Context, userID *uuid.UUID) (*model.Player, error) {
	if userID == nil {
		return nil, nil
	}

	players, err := s.playerUC.GetByUsers(ctx, []uuid.UUID{*userID})
	if err != nil {
		return nil, err
	}

	if len(players) <= 0 {
		return nil, nil
	}

	return &players[0], nil
}

func (s *DefaultService) findFromCookie(request *http.Request, gameID uuid.UUID) (*model.Player, error) {
	value, err := s.cookie.Get(request, cookieName(gameID))
	if errors.Is(err, http.ErrNoCookie) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	playerID, err := uuid.Parse(value)
	if err != nil {
		return nil, err
	}

	players, err := s.playerUC.Get(request.Context(), []uuid.UUID{playerID})
	if err != nil {
		return nil, err
	}

	if len(players) <= 0 {
		return nil, nil
	}

	return &players[0], nil
}

func (s *DefaultService) setPlayerID(writer http.ResponseWriter, id uuid.UUID, gameID uuid.UUID) error {
	return s.cookie.Set(writer, cookieName(gameID), id.String(), 24*time.Hour)
}

func cookieName(gameID uuid.UUID) string {
	return fmt.Sprintf("%s-%s", cookiePlayerID, gameID.String())
}
