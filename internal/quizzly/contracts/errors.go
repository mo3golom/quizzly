package contracts

import "errors"

var (
	ErrQuestionQueueIsEmpty     = errors.New("question queue is empty")
	ErrUnansweredQuestionExists = errors.New("you have unanswered question")
)
