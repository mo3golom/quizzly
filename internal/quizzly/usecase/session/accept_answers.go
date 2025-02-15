package session

import (
	"context"
	"errors"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/pkg/structs"
	"time"

	"github.com/google/uuid"
)

func (u *Usecase) AcceptAnswers(ctx context.Context, in *contracts.AcceptAnswersIn) (*contracts.AcceptAnswersOut, error) {
	var result *contracts.AcceptAnswersOut
	return result, u.trm.Do(ctx, func(ctx context.Context) error {
		if _, err := u.getActiveGame(ctx, in.GameID); err != nil {
			return err
		}

		specificSession, err := u.getSession(ctx, in.PlayerID, in.GameID)
		if err != nil {
			return err
		}
		if specificSession.Status != model.SessionStatusStarted {
			return contracts.ErrNotActiveSessionNotFound
		}

		specificQuestions, err := u.games.GetQuestionsBySpec(ctx, &game.QuestionsSpec{
			IDs: []uuid.UUID{in.QuestionID},
		})
		if err != nil {
			return err
		}
		if len(specificQuestions) == 0 {
			return errors.New("question not found")
		}

		result, err = u.acceptAnswers(&specificQuestions[0], in.Answers)
		result.RightAnswers = specificQuestions[0].GetCorrectAnswers()
		if err != nil {
			return err
		}

		return u.sessions.InsertSessionItem(
			ctx,
			&model.SessionItem{
				SessionID:  specificSession.ID,
				QuestionID: in.QuestionID,
				IsCorrect:  structs.Pointer(result.IsCorrect),
				Answers:    in.Answers,
				AnsweredAt: structs.Pointer(time.Now()),
			},
		)
	})
}

func (u *Usecase) acceptAnswers(question *model.Question, answers []string) (*contracts.AcceptAnswersOut, error) {
	if len(answers) == 0 {
		return nil, errors.New("answers are empty")
	}

	if acceptor, ok := u.optionIDAcceptors[question.Type]; ok {
		return acceptor.Accept(question, answers)
	}

	return nil, errors.New("question type is not supported")
}
