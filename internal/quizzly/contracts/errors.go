package contracts

import "errors"

var (
	ErrQuestionQueueIsEmpty     = errors.New("question queue is empty")
	ErrNotActiveSessionNotFound = errors.New("player's active session not found")
	ErrSessionNotFound          = errors.New("player's session not found")
	ErrSessionNotFinished       = errors.New("player's session not finished")
	ErrGameNotFound             = errors.New("game not found")
	ErrEmptyQuestions           = errors.New("empty questions")
	ErrEmptyAnswerOptions       = errors.New("empty answer options")
)
