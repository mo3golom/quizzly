package contracts

import "errors"

var (
	ErrQuestionQueueIsEmpty     = errors.New("question queue is empty")
	ErrNotActiveSessionNotFound = errors.New("player's active session not found")
	ErrSessionNotFound          = errors.New("player's session not found")
)
