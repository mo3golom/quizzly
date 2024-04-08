package model

import "github.com/google/uuid"

const (
	QuestionTypChoice         QuestionType = "choice"
	QuestionTypOneOfChoice    QuestionType = "one_of_choice"
	QuestionTypMultipleChoice QuestionType = "multiple_choice"
)

type (
	QuestionType string

	Question struct {
		ID            uuid.UUID
		Type          QuestionType
		AnswerOptions []AnswerOption
	}

	AnswerOption struct {
		ID        uuid.UUID
		Answer    string
		IsCorrect bool
	}
)
