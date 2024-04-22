package session

import (
	"errors"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
)

func checkAnswers(question *model.Question, answers []string) (*contracts.AcceptAnswersOut, error) {
	if len(answers) == 0 {
		return nil, errors.New("answers are empty")
	}

	result := &contracts.AcceptAnswersOut{
		IsCorrect: true,
		Details:   make([]contracts.AnswerResult, 0, len(answers)),
	}

	correctAnswers := question.GetCorrectAnswers()
	correctAnswersMap := make(map[string]struct{}, len(correctAnswers))
	for _, answer := range correctAnswers {
		correctAnswersMap[answer.Answer] = struct{}{}
	}

	switch question.Type {
	case model.QuestionTypeChoice:
		if len(answers) > 1 {
			return nil, errors.New("simple choice can't have multiple answers")
		}

		_, ok := correctAnswersMap[answers[0]]
		result.IsCorrect = ok
		result.Details = append(result.Details, contracts.AnswerResult{
			Answer:    answers[0],
			IsCorrect: ok,
		})
	case model.QuestionTypeOneOfChoice:
		for _, answer := range answers {
			_, ok := correctAnswersMap[answer]

			if !ok && result.IsCorrect {
				result.IsCorrect = ok
			}

			result.Details = append(result.Details, contracts.AnswerResult{
				Answer:    answer,
				IsCorrect: ok,
			})
		}
	case model.QuestionTypeMultipleChoice:
		findCount := 0
		for _, answer := range answers {
			_, ok := correctAnswersMap[answer]

			if !ok && result.IsCorrect {
				result.IsCorrect = ok
			}

			if ok {
				findCount++
				delete(correctAnswersMap, answer)
			}

			result.Details = append(result.Details, contracts.AnswerResult{
				Answer:    answer,
				IsCorrect: ok,
			})
		}

		result.IsCorrect = findCount == len(correctAnswers)
	}

	return result, nil
}
