package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type (
	tokenService struct {
		secretKey []byte
	}
)

func newTokenService(secretKey string) *tokenService {
	return &tokenService{
		secretKey: []byte(secretKey),
	}
}

func (s *tokenService) createToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID.String(),
			"exp":     time.Now().Add(ttl).Unix(),
		},
	)

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *tokenService) verifyToken(token string) error {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return fmt.Errorf("invalid token")
	}

	expiredAt, err := jwtToken.Claims.GetExpirationTime()
	if err != nil {
		return err
	}
	if expiredAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	return nil
}

func (s *tokenService) getUserID(token string) (uuid.UUID, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("failed to get user id")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user id: %w", err)
	}

	return userID, nil
}
