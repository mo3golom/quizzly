package model

import "github.com/google/uuid"

const (
	GameStatusCreated  GameStatus = "created"
	GameStatusStarted  GameStatus = "started"
	GameStatusFinished GameStatus = "finished"
)

const (
	GameTypeAsync     GameType = "async"
	GameTypeOnline    GameType = "online"
	GameTypeLocalQuiz GameType = "local_quiz"
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
		LoginRequired    bool
		ShuffleQuestions bool
		ShuffleAnswers   bool
	}

	PlayerGameStatus string

	PlayerGame struct {
		ID     uuid.UUID
		GameID uuid.UUID
		Status PlayerGameStatus
	}
)
