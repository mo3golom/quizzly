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
		ID       uuid.UUID
		Status   GameStatus
		Type     GameType
		Settings GameSettings
	}

	GameSettings struct {
		IsPrivate        bool
		ShuffleQuestions bool
		ShuffleAnswers   bool
	}

	SessionStatus string

	Session struct {
		ID       int64
		PlayerID uuid.UUID
		GameID   uuid.UUID
		Status   SessionStatus
	}

	SessionItem struct {
		ID         int64
		SessionID  int64
		QuestionID uuid.UUID
		Answers    []string
		IsCorrect  *bool
		AnsweredAt *time.Time
	}
)
