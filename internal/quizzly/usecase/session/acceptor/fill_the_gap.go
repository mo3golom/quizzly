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
	if len(correctAnswers) == 0 {
		return nil, errors.New("no correct answers defined")
	}

	distance := levenshteinDistance(correctAnswers[0].Answer, answers[0])
	// Consider the answer correct if the Levenshtein distance is less than 2
	// You can adjust this threshold based on your needs
	isCorrect := distance <= 2

	return &contracts.AcceptAnswersOut{
		IsCorrect: isCorrect,
		Details: []contracts.AnswerResult{
			{
				Answer:    answers[0],
				IsCorrect: isCorrect,
			},
		},
	}, nil
}

func levenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}
	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	minResult := a
	if b < minResult {
		minResult = b
	}
	if c < minResult {
		return c
	}
	return minResult
}
