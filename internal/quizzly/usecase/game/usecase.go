package game

import (
	"context"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type Usecase struct {
	games     game.Repository
	questions question.Repository
	sessions  session.Repository
	template  transactional.Template
}

func NewUsecase(
	games game.Repository,
	questions question.Repository,
	sessions session.Repository,
	template transactional.Template,
) contracts.GameUsecase {
	return &Usecase{
		games:     games,
		questions: questions,
		sessions:  sessions,
		template:  template,
	}
}

func (u *Usecase) Create(ctx context.Context, in *contracts.CreateGameIn) (uuid.UUID, error) {
	id := uuid.New()

	return id, u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.games.Insert(
			ctx,
			tx,
			&model.Game{
				ID:       id,
				AuthorID: in.AuthorID,
				Status:   model.GameStatusCreated,
				Type:     model.GameTypeAsync,
				Settings: in.Settings,
			},
		)
	})

}

func (u *Usecase) Start(ctx context.Context, id uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		specificGame, err := u.games.GetWithTx(ctx, tx, id)
		if err != nil {
			return err
		}

		specificGame.Status = model.GameStatusStarted
		return u.games.Update(ctx, tx, specificGame)
	})
}

func (u *Usecase) Finish(ctx context.Context, id uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		specificGame, err := u.games.GetWithTx(ctx, tx, id)
		if err != nil {
			return err
		}

		specificGame.Status = model.GameStatusFinished
		return u.games.Update(ctx, tx, specificGame)
	})
}

func (u *Usecase) Get(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	var specificGame *model.Game
	return specificGame, u.template.Execute(ctx, func(tx transactional.Tx) error {
		var err error
		specificGame, err = u.games.GetWithTx(ctx, tx, id)
		return err
	})
}

func (u *Usecase) GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]model.Game, error) {
	return u.games.GetByAuthorID(ctx, authorID)
}

func (u *Usecase) AddQuestion(ctx context.Context, gameID uuid.UUID, questionID ...uuid.UUID) error {
	if len(questionID) == 0 {
		return nil
	}

	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		specificGame, err := u.games.GetWithTx(ctx, tx, gameID)
		if err != nil {
			return err
		}

		specificQuestions, err := u.questions.GetBySpec(ctx, &question.Spec{
			IDs: questionID,
		})
		if err != nil {
			return err
		}

		return u.games.InsertGameQuestions(
			ctx,
			tx,
			specificGame.ID,
			slices.SafeMap(specificQuestions, func(question model.Question) uuid.UUID {
				return question.ID
			}),
		)
	})
}

func (u *Usecase) GetQuestions(ctx context.Context, gameID uuid.UUID) ([]uuid.UUID, error) {
	var result []uuid.UUID
	return result, u.template.Execute(ctx, func(tx transactional.Tx) error {
		var err error
		result, err = u.games.GetQuestionIDsBySpec(ctx, tx, &game.Spec{
			GameID: gameID,
		})
		return err
	})
}

func (u *Usecase) GetStatistics(ctx context.Context, id uuid.UUID) (*model.GameStatistics, error) {
	var questionsCount int64
	var participantsCount int64
	err := u.template.Execute(ctx, func(tx transactional.Tx) error {
		questions, err := u.games.GetQuestionIDsBySpec(ctx, tx, &game.Spec{
			GameID: id,
		})
		if err != nil {
			return err
		}

		questionsCount = int64(len(questions))
		return nil
	})

	sessions, err := u.sessions.GetSessionsByGameID(ctx, id)
	if err != nil {
		return nil, err
	}

	participantsCount = int64(len(sessions))

	return &model.GameStatistics{
		QuestionsCount:    questionsCount,
		ParticipantsCount: participantsCount,
		CompletionRate:    calculateCompletionRate(sessions),
	}, nil
}

func calculateCompletionRate(sessions []model.SessionExtended) int64 {
	if len(sessions) == 0 {
		return 0
	}

	var sum int64
	var count int64
	for _, item := range sessions {
		sum += item.CompletionRate()
		count++
	}

	return sum / count
}
