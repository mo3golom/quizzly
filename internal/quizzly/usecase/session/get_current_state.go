package session

import (
	"context"
	cryptoRand "crypto/rand"
	"errors"
	"github.com/google/uuid"
	"math/big"
	"math/rand"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"
	"sort"
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

		gameQuestions, err := u.games.GetQuestionIDsBySpec(
			ctx,
			tx,
			&game.QuestionSpec{
				GameID: gameID,
			},
		)
		if err != nil {
			return err
		}

		sessionItems, err := u.sessions.GetSessionBySpecWithTx(
			ctx,
			tx,
			&session.ItemSpec{
				PlayerID: playerID,
				GameID:   gameID,
			},
		)
		if err != nil {
			return err
		}

		specificUnansweredQuestion, err := getUnansweredQuestionID(specificGame, gameQuestions, sessionItems)
		if err != nil {
			return err
		}

		specificQuestions, err := u.questions.GetBySpec(ctx, &question.Spec{
			IDs: []uuid.UUID{specificUnansweredQuestion.ID},
		})
		if err != nil {
			return err
		}
		if len(specificQuestions.Result) == 0 {
			return errors.New("question not found")
		}

		currentQuestion := specificQuestions.Result[0]
		if specificGame.Settings.ShuffleAnswers {
			tempAnswerOptions := currentQuestion.AnswerOptions
			rand.Shuffle(len(tempAnswerOptions), func(i, j int) {
				tempAnswerOptions[i], tempAnswerOptions[j] = tempAnswerOptions[j], tempAnswerOptions[i]
			})
			currentQuestion.AnswerOptions = tempAnswerOptions
		}

		result = &contracts.SessionState{
			CurrentQuestion: &currentQuestion,
			Progress: contracts.Progress{
				Total: int64(len(gameQuestions)),
				Answered: int64(len(
					slices.Filter(sessionItems, func(item model.SessionItem) bool {
						return item.AnsweredAt != nil
					}),
				)),
			},
		}

		if !specificUnansweredQuestion.IsNew {
			return nil
		}
		return u.sessions.InsertSessionItem(ctx, tx, &model.SessionItem{
			SessionID:  specificSession.ID,
			QuestionID: currentQuestion.ID,
		})
	})
}

func getUnansweredQuestionID(specificGame *model.Game, gameQuestions []game.GameQuestion, sessionItems []model.SessionItem) (*unansweredQuestion, error) {
	unansweredSessionItems := slices.Filter(sessionItems, func(item model.SessionItem) bool {
		return item.AnsweredAt == nil
	})
	if len(unansweredSessionItems) > 0 {
		return &unansweredQuestion{
			ID:    unansweredSessionItems[0].QuestionID,
			IsNew: false,
		}, nil
	}

	answeredSessionItemsMap := make(map[uuid.UUID]struct{}, len(sessionItems))
	for _, item := range sessionItems {
		answeredSessionItemsMap[item.QuestionID] = struct{}{}
	}

	unansweredQuestions := slices.Filter(gameQuestions, func(question game.GameQuestion) bool {
		_, ok := answeredSessionItemsMap[question.ID]
		return !ok
	})
	if len(unansweredQuestions) == 0 {
		return nil, contracts.ErrQuestionQueueIsEmpty
	}

	sort.Slice(unansweredQuestions, func(i, j int) bool {
		return unansweredQuestions[i].Sort > unansweredQuestions[j].Sort
	})

	questionIndex := 0
	if specificGame.Settings.ShuffleQuestions {
		randomNumber, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(unansweredQuestions))))
		questionIndex = int(randomNumber.Int64())
	}

	return &unansweredQuestion{
		ID:    unansweredQuestions[questionIndex].ID,
		IsNew: true,
	}, nil
}
