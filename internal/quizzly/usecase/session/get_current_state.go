package session

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"
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
			&game.Spec{
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
