package session

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs"
	"quizzly/pkg/transactional"
	"time"
)

func (u *Usecase) AcceptAnswers(ctx context.Context, in *contracts.AcceptAnswersIn) (*contracts.AcceptAnswersOut, error) {
	var result *contracts.AcceptAnswersOut
	return result, u.template.Execute(ctx, func(tx transactional.Tx) error {
		if _, err := u.getActiveGame(ctx, tx, in.GameID); err != nil {
			return err
		}

		specificSession, err := u.getSession(ctx, tx, in.PlayerID, in.GameID)
		if err != nil {
			return err
		}
		if specificSession.Status != model.SessionStatusStarted {
			return contracts.ErrNotActiveSessionNotFound
		}

		specificPlayerSessionItems, err := u.sessions.GetSessionBySpecWithTx(ctx, tx, &session.ItemSpec{
			PlayerID:   in.PlayerID,
			GameID:     in.GameID,
			QuestionID: &in.QuestionID,
		})
		if err != nil {
			return err
		}
		if len(specificPlayerSessionItems) == 0 {
			return errors.New("player session is empty")
		}

		specificPlayerSessionItem := specificPlayerSessionItems[0]
		if specificPlayerSessionItem.AnsweredAt != nil {
			return errors.New("question is already answered")
		}

		specificQuestions, err := u.questions.GetBySpec(ctx, &question.Spec{
			IDs: []uuid.UUID{in.QuestionID},
		})
		if err != nil {
			return err
		}
		if len(specificQuestions) == 0 {
			return errors.New("question not found")
		}

		result, err = u.checkAnswers(&specificQuestions[0], in.Answers)
		if err != nil {
			return err
		}

		specificPlayerSessionItem.IsCorrect = structs.Pointer(result.IsCorrect)
		specificPlayerSessionItem.Answers = in.Answers
		specificPlayerSessionItem.AnsweredAt = structs.Pointer(time.Now())
		return u.sessions.UpdateSessionItem(
			ctx,
			tx,
			&specificPlayerSessionItem,
		)
	})
}

func (u *Usecase) checkAnswers(question *model.Question, answers []model.AnswerOptionID) (*contracts.AcceptAnswersOut, error) {
	if len(answers) == 0 {
		return nil, errors.New("answers are empty")
	}

	checker, ok := u.checkers[question.Type]
	if !ok {
		return nil, errors.New("question type is not supported")
	}

	return checker.Check(question, answers)
}
