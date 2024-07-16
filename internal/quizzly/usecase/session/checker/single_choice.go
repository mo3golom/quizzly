package checker

import (
	"errors"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
)

type SingleChoiceChecker struct{}

func NewSingleChoiceChecker() *SingleChoiceChecker {
	return &SingleChoiceChecker{}
}

func (c *SingleChoiceChecker) Check(question *model.Question, answers []model.AnswerOptionID) (*contracts.AcceptAnswersOut, error) {
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
				Answer:    answers[0],
				IsCorrect: ok,
			},
		},
	}, nil
}