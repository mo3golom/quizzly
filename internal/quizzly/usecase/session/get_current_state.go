package session

import (
	"context"
	"errors"
	"math/rand/v2"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"
	"strconv"

	"github.com/google/uuid"
)

func (u *Usecase) GetCurrentState(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*contracts.SessionState, error) {
	var result *contracts.SessionState
	return result, u.template.Execute(ctx, func(tx transactional.Tx) error {
		specificGame, err := u.getActiveGame(ctx, tx, gameID)
		if err != nil {
			return err
		}

		specificSession, err := u.getSession(ctx, tx, playerID, gameID)
		if err != nil {
			return err
		}
		if specificSession.Status == model.SessionStatusFinished {
			result = &contracts.SessionState{
				Status: specificSession.Status,
			}
			return nil
		}

		questionMap, err := u.getQuestionMap(ctx, tx, gameID)
		if err != nil {
			return err
		}

		sessionItems, err := u.sessions.GetSessionBySpecWithTx(ctx, tx, &session.ItemSpec{
			PlayerID: playerID,
			GameID:   gameID,
		})
		if err != nil {
			return err
		}

		currentQuestion, err := findUnansweredQuestion(specificGame, questionMap, sessionItems)
		if err != nil {
			return err
		}

		if specificGame.Settings.ShuffleAnswers {
			tempAnswerOptions := currentQuestion.AnswerOptions
			rand.Shuffle(len(tempAnswerOptions), func(i, j int) {
				tempAnswerOptions[i], tempAnswerOptions[j] = tempAnswerOptions[j], tempAnswerOptions[i]
			})
			currentQuestion.AnswerOptions = tempAnswerOptions
		}

		result = &contracts.SessionState{
			CurrentQuestion: currentQuestion,
			Progress: contracts.Progress{
				Total: int64(questionMap.Len()),
				Answered: int64(len(
					slices.Filter(sessionItems, func(item model.SessionItem) bool {
						return item.AnsweredAt != nil
					}),
				)),
			},
		}
		return nil
	})
}

func findUnansweredQuestion(specificGame *model.Game, questionMap *model.QuestionMap, sessionItems []model.SessionItem) (*model.Question, error) {
	if len(sessionItems) == 0 {
		if questionMap.Empty() {
			return nil, contracts.ErrEmptyQuestions
		}

		if specificGame.Settings.ShuffleQuestions {
			return questionMap.GetRandomQuestion(), nil
		}
		return questionMap.GetFirst(), nil
	}

	answeredMap := make(map[uuid.UUID]bool, len(sessionItems))
	for _, item := range sessionItems {
		answeredMap[item.QuestionID] = true
	}

	var latestSessionItem model.SessionItem
	for _, item := range sessionItems {
		if latestSessionItem.AnsweredAt == nil ||
			(item.AnsweredAt != nil && item.AnsweredAt.After(*latestSessionItem.AnsweredAt)) {
			latestSessionItem = item
		}
	}

	lastQuestion, err := questionMap.GetQuestion(latestSessionItem.QuestionID)
	if err != nil {
		return nil, err
	}

	if specificGame.Settings.ShuffleQuestions {
		var unansweredQuestions []uuid.UUID
		for _, id := range questionMap.GetIDs() {
			id := id
			if answeredMap[id] {
				continue
			}

			unansweredQuestions = append(unansweredQuestions, id)
		}

		if len(unansweredQuestions) == 0 {
			return nil, contracts.ErrQuestionQueueIsEmpty
		}

		question, err := questionMap.GetQuestion(unansweredQuestions[rand.IntN(len(unansweredQuestions))])
		if errors.Is(err, model.ErrQuestionNotFound) {
			return nil, contracts.ErrQuestionQueueIsEmpty
		}
		if err != nil {
			return nil, err
		}

		return question, nil
	}

	var nextQuestion *model.Question
	if len(latestSessionItem.Answers) > 0 {
		selectedAnswer, err := strconv.Atoi(latestSessionItem.Answers[0])
		if err != nil {
			return nil, err
		}

		nextQuestion, err = questionMap.GetNextQuestion(lastQuestion.ID, model.AnswerOptionID(selectedAnswer))
		if err != nil && !errors.Is(err, model.ErrQuestionNotFound) {
			return nil, err
		}
	}

	if nextQuestion != nil {
		return nextQuestion, nil
	}

	for i, id := range questionMap.GetIDs() {
		if id != lastQuestion.ID || i >= questionMap.Len()-1 {
			continue
		}

		question, err := questionMap.GetQuestion(questionMap.GetIDs()[i+1])
		if errors.Is(err, model.ErrQuestionNotFound) {
			return nil, contracts.ErrQuestionQueueIsEmpty
		}
		if err != nil {
			return nil, err
		}

		return question, nil
	}

	return nil, contracts.ErrQuestionQueueIsEmpty
}

func (u *Usecase) getQuestionMap(ctx context.Context, tx transactional.Tx, gameID uuid.UUID) (*model.QuestionMap, error) {
	questions, err := u.games.GetQuestionsBySpecWithTx(ctx, tx, &game.QuestionsSpec{
		GameID: &gameID,
		Order: &game.Order{
			Field:     "created_at",
			Direction: "asc",
		},
	})
	if err != nil {
		return nil, err
	}

	return model.NewQuestionMap(questions), nil
}
