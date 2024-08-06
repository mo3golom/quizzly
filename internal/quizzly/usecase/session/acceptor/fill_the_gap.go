package acceptor

import (
	"errors"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"strings"
)

type FillTheGapAcceptor struct{}

func NewFillTheGapAcceptor() *FillTheGapAcceptor {
	return &FillTheGapAcceptor{}
}

func (a *FillTheGapAcceptor) Accept(question *model.Question, answers []string) (*contracts.AcceptAnswersOut, error) {
	if len(answers) > 1 {
		return nil, errors.New("fill the gap can't have multiple answers")
	}

	correctAnswers := question.GetCorrectAnswers()
	correctAnswersMap := make(map[string]struct{}, len(correctAnswers))
	for _, answer := range correctAnswers {
		correctAnswersMap[strings.ToLower(answer.Answer)] = struct{}{}
	}

	_, ok := correctAnswersMap[strings.ToLower(answers[0])]
	return &contracts.AcceptAnswersOut{
		IsCorrect: ok,
		Details: []contracts.AnswerResult{
			{
				Answer:    answers[0],
				IsCorrect: ok,
			},
		},
	}, nil
}
