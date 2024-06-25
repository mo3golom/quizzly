package handlers

import (
	"github.com/google/uuid"
	"quizzly/internal/quizzly/model"
)

var (
	QuestionTypeColors = map[model.QuestionType]QuestionColor{
		model.QuestionTypeChoice: {
			Color: "blue",
			AnswerOptionColors: []string{
				"orange",
				"pink",
				"amber",
				"red",
			},
		},
		model.QuestionTypeMultipleChoice: {
			Color: "amber",
			AnswerOptionColors: []string{
				"indigo",
				"pink",
				"blue",
				"red",
			},
		},
	}
)

type (
	QuestionColor struct {
		Color              string
		AnswerOptionColors []string
	}

	Game struct {
		ID     uuid.UUID
		Status model.GameStatus
		Link   string
	}

	GameStatistics struct {
		QuestionsCount    int
		ParticipantsCount int
		CompletionRate    int
	}

	Session struct {
		PlayerID        uuid.UUID
		CurrentQuestion *Question
		AnswerResult    *AnswerResult
		Progress        SessionProgress
		Status          model.SessionStatus
		Statistics      *SessionStatistics
	}

	Participant struct {
		PlayerID       uuid.UUID
		SessionStatus  model.SessionStatus
		CompletionRate int
	}

	Question struct {
		ID            uuid.UUID
		Text          string
		Type          model.QuestionType
		AnswerOptions []AnswerOption
		Color         string
	}

	AnswerOption struct {
		ID    uuid.UUID
		Text  string
		Color string
	}

	AnswerResult struct {
		IsCorrect bool
	}

	SessionProgress struct {
		Total    int
		Answered int
	}

	SessionStatistics struct {
		QuestionsCount      int
		CorrectAnswersCount int
	}
)
