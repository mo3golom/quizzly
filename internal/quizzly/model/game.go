package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	GameStatusCreated  GameStatus = "created"
	GameStatusStarted  GameStatus = "started"
	GameStatusFinished GameStatus = "finished"
)

const (
	GameTypeAsync GameType = "async"
)

const (
	SessionStatusStarted  SessionStatus = "started"
	SessionStatusFinished SessionStatus = "finished"
)

type (
	GameStatus string
	GameType   string

	Game struct {
		ID        uuid.UUID
		AuthorID  uuid.UUID
		Status    GameStatus
		Type      GameType
		Title     *string
		Settings  GameSettings
		CreatedAt time.Time
	}

	GameSettings struct {
		IsPrivate        bool
		ShuffleQuestions bool
		ShuffleAnswers   bool
	}

	GameStatistics struct {
		QuestionsCount    int64
		ParticipantsCount int64
		CompletionRate    int64
	}

	SessionStatus string

	Session struct {
		ID       int64
		PlayerID uuid.UUID
		GameID   uuid.UUID
		Status   SessionStatus
	}

	SessionExtended struct {
		Session
		Items []SessionItem
	}

	SessionItem struct {
		ID         int64
		SessionID  int64
		QuestionID uuid.UUID
		Answers    []AnswerOptionID
		IsCorrect  *bool
		AnsweredAt *time.Time
		CreatedAt  time.Time
	}

	SessionStatistics struct {
		QuestionsCount      int64
		CorrectAnswersCount int64
	}
)

func (s *SessionExtended) CompletionRate() int64 {
	if len(s.Items) == 0 {
		return 0
	}

	total := int64(len(s.Items))
	correctAnswers := int64(0)
	for _, item := range s.Items {
		if item.IsCorrect == nil || !*item.IsCorrect {
			continue
		}

		correctAnswers++
	}

	return (correctAnswers * 100) / total
}
