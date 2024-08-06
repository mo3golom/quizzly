package acceptor

import (
	"errors"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"strconv"
)

type SingleChoiceAcceptor struct{}

func NewSingleChoiceAcceptor() *SingleChoiceAcceptor {
	return &SingleChoiceAcceptor{}
}

func (a *SingleChoiceAcceptor) Accept(question *model.Question, answers []model.AnswerOptionID) (*contracts.AcceptAnswersOut, error) {
	if len(answers) > 1 {
		return nil, errors.New("simple choice can't have multiple answers")
	}

	correctAnswers := question.GetCorrectAnswers()
	correctAnswersMap := make(map[model.AnswerOptionID]struct{}, len(correctAnswers))
	for _, answer := range correctAnswers {
		correctAnswersMap[answer.ID] = struct{}{}
	}

	_, ok := correctAnswersMap[answers[0]]
	return &contracts.AcceptAnswersOut{
		IsCorrect: ok,
		Details: []contracts.AnswerResult{
			{
				Answer:    strconv.Itoa(int(answers[0])),
				IsCorrect: ok,
			},
		},
	}, nil
}
