package handlers

import (
	"quizzly/internal/quizzly/model"
	"time"

	"github.com/google/uuid"
)

var (
	QuestionTypePublicColors = QuestionColor{
		Color: "blue",
		AnswerOptionColors: []string{
			"orange",
			"pink",
			"amber",
			"red",
		},
	}
)

type (
	QuestionColor struct {
		Color              string
		AnswerOptionColors []string
	}

	Game struct {
		ID        uuid.UUID
		Status    model.GameStatus
		Title     string
		CreatedAt time.Time
		Settings  GameSettings
	}

	GameSettings struct {
		ShuffleQuestions bool
		ShuffleAnswers   bool
		ShowRightAnswers bool
		InputCustomName  bool
		IsPrivate        bool
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

	SessionItem struct {
		QuestionText  string
		QuestionImage *string
		Answers       []SessionItemAnswer
	}

	SessionItemAnswer struct {
		AnswerText       string
		IsCorrect        bool
		IsPlayerAnswered bool
	}

	Participant struct {
		PlayerID       uuid.UUID
		SessionStatus  model.SessionStatus
		CompletionRate int
	}

	Question struct {
		ID            uuid.UUID
		ImageID       *string
		Text          string
		Type          model.QuestionType
		AnswerOptions []AnswerOption
		Color         string
	}

	AnswerOption struct {
		ID    int64
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

	SessionItemStatistics struct {
		PlayerName                    string
		CompletionRate                int
		SessionStatus                 model.SessionStatus
		SessionStartedAt              *time.Time
		SessionLastQuestionAnsweredAt *time.Time
	}
)
