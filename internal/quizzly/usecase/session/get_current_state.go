package session

import (
	"context"
	"github.com/google/uuid"
	"math/rand/v2"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs/collections/slices"
)

func (u *Usecase) GetCurrentState(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (*contracts.SessionState, error) {
	var result *contracts.SessionState
	return result, u.trm.Do(ctx, func(ctx context.Context) error {
		specificGame, err := u.getActiveGame(ctx, gameID)
		if err != nil {
			return err
		}

		specificSession, err := u.getSession(ctx, playerID, gameID)
		if err != nil {
			return err
		}
		if specificSession.Status == model.SessionStatusFinished {
			result = &contracts.SessionState{
				Status: specificSession.Status,
			}
			return nil
		}

		questionList, err := u.games.GetQuestionsBySpec(ctx, &game.QuestionsSpec{
			GameID: &gameID,
		})
		if err != nil {
			return err
		}

		sessionItems, err := u.sessions.GetSessionBySpec(ctx, &session.ItemSpec{
			PlayerID: playerID,
			GameID:   gameID,
		})
		if err != nil {
			return err
		}

		currentQuestion, err := findUnansweredQuestion(specificGame, questionList, sessionItems)
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
				Total: int64(len(questionList)),
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

func findUnansweredQuestion(specificGame *model.Game, questions []model.Question, sessionItems []model.SessionItem) (*model.Question, error) {
	if len(questions) == 0 {
		return nil, contracts.ErrEmptyQuestions
	}

	if len(sessionItems) == 0 {
		index := 0
		if specificGame.Settings.ShuffleQuestions {
			index = rand.IntN(len(questions))
		}

		return &questions[index], nil
	}

	var latestSessionItem model.SessionItem
	answeredMap := make(map[uuid.UUID]bool, len(sessionItems))
	for _, item := range sessionItems {
		answeredMap[item.QuestionID] = true

		if latestSessionItem.AnsweredAt == nil ||
			(item.AnsweredAt != nil && item.AnsweredAt.After(*latestSessionItem.AnsweredAt)) {
			latestSessionItem = item
		}
	}

	questionIndexMap := make(map[uuid.UUID]int, len(questions))
	unansweredQuestions := make([]uuid.UUID, 0, len(questions))
	for i, item := range questions {
		questionIndexMap[item.ID] = i

		if answeredMap[item.ID] {
			continue
		}
		unansweredQuestions = append(unansweredQuestions, item.ID)
	}

	if len(unansweredQuestions) == 0 {
		return nil, contracts.ErrQuestionQueueIsEmpty
	}

	if specificGame.Settings.ShuffleQuestions {
		nextQuestionIndex, ok := questionIndexMap[unansweredQuestions[rand.IntN(len(unansweredQuestions))]]
		if !ok {
			return nil, contracts.ErrQuestionQueueIsEmpty
		}

		return &questions[nextQuestionIndex], nil
	}

	lastQuestionIndex, ok := questionIndexMap[latestSessionItem.QuestionID]
	if !ok || lastQuestionIndex >= len(questions)-1 {
		return nil, contracts.ErrQuestionQueueIsEmpty
	}

	nextQuestion := questions[lastQuestionIndex+1]
	return &nextQuestion, nil
}
