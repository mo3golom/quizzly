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
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"
	"strconv"
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
		if len(specificQuestions.Result) == 0 {
			return errors.New("question not found")
		}

		result, err = u.acceptAnswers(&specificQuestions.Result[0], in.Answers)
		result.RightAnswers = specificQuestions.Result[0].GetCorrectAnswers()
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

func (u *Usecase) acceptAnswers(question *model.Question, answers []string) (*contracts.AcceptAnswersOut, error) {
	if len(answers) == 0 {
		return nil, errors.New("answers are empty")
	}

	if acceptor, ok := u.textAcceptors[question.Type]; ok {
		return acceptor.Accept(question, answers)
	}

	if acceptor, ok := u.optionIDAcceptors[question.Type]; ok {
		convertedAnswers, err := slices.Map(answers, func(i string) (model.AnswerOptionID, error) {
			id, err := strconv.Atoi(i)
			return model.AnswerOptionID(id), err
		})
		if err != nil {
			return nil, err
		}

		return acceptor.Accept(question, convertedAnswers)
	}

	return nil, errors.New("question type is not supported")
}
