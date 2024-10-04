package player

import (
	"net/http"
	"quizzly/internal/quizzly/model"
)

type (
	Service interface {
		GetPlayer(writer http.ResponseWriter, request *http.Request, customName ...string) (*model.Player, error)
	}
)
