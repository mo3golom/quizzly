package acceptor

import (
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"strconv"
)

type OneOfChoiceAcceptor struct{}

func NewOneOfChoiceAcceptor() *OneOfChoiceAcceptor {
	return &OneOfChoiceAcceptor{}
}

func (a *OneOfChoiceAcceptor) Accept(question *model.Question, answers []string) (*contracts.AcceptAnswersOut, error) {
	correctAnswers := question.GetCorrectAnswers()
	correctAnswersMap := make(map[string]struct{}, len(correctAnswers))
	for _, answer := range correctAnswers {
		correctAnswersMap[strconv.FormatInt(int64(answer.ID), 10)] = struct{}{}
	}

	result := &contracts.AcceptAnswersOut{
		IsCorrect: true,
		Details:   make([]contracts.AnswerResult, 0, len(answers)),
	}
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

	return result, nil
}
