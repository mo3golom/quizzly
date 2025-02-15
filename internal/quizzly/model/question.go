package model

import (
	"github.com/google/uuid"
	"time"
)

const (
	QuestionTypeChoice         QuestionType = "choice"
	QuestionTypeOneOfChoice    QuestionType = "one_of_choice"   // Может быть выбран любой правильный вариант ответа
	QuestionTypeMultipleChoice QuestionType = "multiple_choice" // Должны быть выбраны все правильные ответы
	QuestionTypeFillTheGap     QuestionType = "fill_the_gap"    // Нужно ввести правильный ответ
)

type (
	QuestionType   string
	AnswerOptionID int64

	Question struct {
		ID            uuid.UUID
		GameID        uuid.UUID
		Text          string
		Type          QuestionType
		ImageID       *string
		AnswerOptions []AnswerOption
		CreatedAt     time.Time
	}

	AnswerOption struct {
		ID        AnswerOptionID
		Answer    string
		IsCorrect bool
	}
)

func (q Question) GetCorrectAnswers() []AnswerOption {
	result := make([]AnswerOption, 0, len(q.AnswerOptions))
	for _, answer := range q.AnswerOptions {
		answer := answer
		if !answer.IsCorrect {
			continue
		}

		result = append(result, answer)
	}

	return result
}
