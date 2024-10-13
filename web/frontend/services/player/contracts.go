package player

import (
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/model"
)

type (
	Service interface {
		GetPlayer(writer http.ResponseWriter, request *http.Request, gameID uuid.UUID, customName ...string) (*model.Player, error)
	}
)
