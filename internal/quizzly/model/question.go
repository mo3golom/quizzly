package model

import "github.com/google/uuid"

const (
	QuestionTypeChoice         QuestionType = "choice"
	QuestionTypeOneOfChoice    QuestionType = "one_of_choice"
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
)

type (
	QuestionType string

	Question struct {
		ID            uuid.UUID
		Text          string
		Type          QuestionType
		AnswerOptions []AnswerOption
	}

	AnswerOption struct {
		Answer    string
		IsCorrect bool
	}
)

func (q Question) AnswersIsCorrect(answers []string) bool {
	if len(answers) == 0 {
		return false
	}

	correctAnswers := q.GetCorrectAnswers()
	correctAnswersMap := make(map[string]struct{}, len(correctAnswers))
	for _, answer := range correctAnswers {
		correctAnswersMap[answer.Answer] = struct{}{}
	}

	switch q.Type {
	case QuestionTypeChoice:
		if len(answers) > 1 {
			return false
		}

		_, ok := correctAnswersMap[answers[0]]
		return ok
	case QuestionTypeOneOfChoice:
		for _, answer := range answers {
			if _, ok := correctAnswersMap[answer]; ok {
				continue
			}

			return false
		}

		return true
	case QuestionTypeMultipleChoice:
		findCount := 0
		for _, answer := range answers {
			if _, ok := correctAnswersMap[answer]; ok {
				findCount++
				delete(correctAnswersMap, answer)

				continue
			}

			return false
		}

		return findCount == len(correctAnswers)
	}

	return false
}

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
