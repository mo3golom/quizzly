package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"quizzly/cmd"
	"quizzly/internal/quizzly"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	if _, err := os.Stat(".env"); err == nil {
		// path/to/whatever exists
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}

	db := cmd.MustInitDB(ctx)
	template := transactional.NewTemplate(db)

	log := cmd.MustInitLogger()
	defer log.Flush()

	quizzlyConfig := quizzly.NewConfiguration(
		db,
		template,
	)

	if err := run(ctx, quizzlyConfig); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, quizzly *quizzly.Configuration) error {
	questionIDs, err := prepareQuestions(ctx, quizzly.Question.MustGet(), nil)
	fmt.Println("questions:", questionIDs)

	gameID, err := prepareGame(ctx, quizzly.Game.MustGet(), err)
	fmt.Println("game:", gameID)
	err = prepareGameQuestions(ctx, quizzly.Game.MustGet(), gameID, questionIDs, err)
	err = startGame(ctx, quizzly.Game.MustGet(), gameID, err)
	playerID, err := startPlayerSession(ctx, quizzly.Session.MustGet(), gameID, err)
	fmt.Println("player:", playerID)

	for {
		if err != nil {
			return err
		}
		question, err := getQuestion(ctx, quizzly.Session.MustGet(), gameID, playerID, err)

		if errors.Is(err, contracts.ErrQuestionQueueIsEmpty) {
			fmt.Println("question queue is empty")
			return quizzly.Session.MustGet().Finish(ctx, gameID, playerID)
		}

		fmt.Println("question:", question)

		answer := question.AnswerOptions[rand.Intn(len(question.AnswerOptions)-1)]

		result, err := answerQuestion(ctx, quizzly.Session.MustGet(), gameID, playerID, question.ID, []string{answer.Answer}, err)
		fmt.Println("answer result:", result)
	}
}

func prepareQuestions(ctx context.Context, question contracts.QuestionUsecase, err error) ([]uuid.UUID, error) {
	if err != nil {
		return nil, err
	}

	questionID := uuid.New()
	err = question.Create(
		ctx,
		&model.Question{
			ID:   questionID,
			Text: "Что говорит Боромир в этом легендарном меме из фильма «Властелин колец»?",
			Type: model.QuestionTypeChoice,
			AnswerOptions: []model.AnswerOption{
				{
					Answer:    "Нельзя просто так взять и зайти в Мордор",
					IsCorrect: true,
				},
				{
					Answer: "Нельзя просто так взять и надеть кольцо",
				},
				{
					Answer: "Нельзя просто так взять и убить Саурона",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return []uuid.UUID{questionID}, nil
}

func prepareGame(ctx context.Context, game contracts.GameUsecase, err error) (uuid.UUID, error) {
	if err != nil {
		return uuid.Nil, err
	}

	gameID, err := game.Create(ctx, &contracts.CreateGameIn{
		Type:     model.GameTypeAsync,
		Settings: model.GameSettings{},
	})
	if err != nil {
		return uuid.Nil, err
	}

	return gameID, nil
}

func prepareGameQuestions(ctx context.Context, game contracts.GameUsecase, gameID uuid.UUID, questionIDs []uuid.UUID, err error) error {
	if err != nil {
		return err
	}

	for _, questionID := range questionIDs {
		err = game.AddQuestion(ctx, gameID, questionID)
		if err != nil {
			return err
		}
	}

	return err
}

func startGame(ctx context.Context, game contracts.GameUsecase, gameID uuid.UUID, err error) error {
	if err != nil {
		return err
	}

	return game.Start(ctx, gameID)
}

func startPlayerSession(ctx context.Context, session contracts.SessionUsecase, gameID uuid.UUID, err error) (uuid.UUID, error) {
	if err != nil {
		return uuid.Nil, err
	}

	playerID := uuid.New()
	err = session.Start(ctx, gameID, playerID)
	if err != nil {
		return uuid.Nil, err
	}

	return playerID, nil
}

func getQuestion(ctx context.Context, session contracts.SessionUsecase, gameID uuid.UUID, playerID uuid.UUID, err error) (*model.Question, error) {
	if err != nil {
		return nil, err
	}

	return session.NextQuestion(ctx, gameID, playerID)
}

func answerQuestion(
	ctx context.Context,
	session contracts.SessionUsecase,
	gameID uuid.UUID,
	playerID uuid.UUID,
	questionID uuid.UUID,
	answers []string,
	err error) (*contracts.AcceptAnswersOut, error) {
	if err != nil {
		return nil, err
	}

	return session.AcceptAnswers(
		ctx,
		&contracts.AcceptAnswersIn{
			GameID:     gameID,
			PlayerID:   playerID,
			QuestionID: questionID,
			Answers:    answers,
		},
	)
}
