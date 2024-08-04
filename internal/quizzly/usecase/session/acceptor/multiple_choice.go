package acceptor

import (
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"strconv"
)

type MultipleChoiceAcceptor struct{}

func NewMultipleChoiceAcceptor() *MultipleChoiceAcceptor {
	return &MultipleChoiceAcceptor{}
}

func (a *MultipleChoiceAcceptor) Accept(question *model.Question, answers []model.AnswerOptionID) (*contracts.AcceptAnswersOut, error) {
	correctAnswers := question.GetCorrectAnswers()
	correctAnswersMap := make(map[model.AnswerOptionID]struct{}, len(correctAnswers))
	for _, answer := range correctAnswers {
		correctAnswersMap[answer.ID] = struct{}{}
	}

	result := &contracts.AcceptAnswersOut{
		IsCorrect: true,
		Details:   make([]contracts.AnswerResult, 0, len(answers)),
	}
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
			Answer:    strconv.Itoa(int(answer)),
			IsCorrect: ok,
		})
	}

	result.IsCorrect = findCount == len(correctAnswers)

	return result, nil
}
