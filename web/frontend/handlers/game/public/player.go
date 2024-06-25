package public

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

func getPlayerID(request *http.Request) uuid.UUID {
	cookie, err := request.Cookie(cookiePlayerID)
	if err != nil {
		return uuid.Nil
	}

	playerID, err := uuid.Parse(cookie.Value)
	if err != nil {
		return uuid.Nil
	}

	return playerID
}

func setPlayerID(writer http.ResponseWriter, id uuid.UUID) {
	cookie := http.Cookie{
		Name:     cookiePlayerID,
		Value:    id.String(),
		Path:     "/",
		Expires:  time.Now().Add(47 * time.Hour),
		MaxAge:   172800,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(writer, &cookie)
}
