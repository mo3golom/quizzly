package game

import (
	"context"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/question"
	"quizzly/pkg/transactional"

	"github.com/google/uuid"
)

type Usecase struct {
	games     game.Repository
	questions question.Repository
	template  transactional.Template
}

func NewUsecase(
	games game.Repository,
	questions question.Repository,
	template transactional.Template,
) contracts.GameUsecase {
	return &Usecase{
		games:     games,
		questions: questions,
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

func (u *Usecase) AddQuestion(ctx context.Context, gameID uuid.UUID, questionID uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		specificGame, err := u.games.GetWithTx(ctx, tx, gameID)
		if err != nil {
			return err
		}

		specificQuestion, err := u.questions.Get(ctx, questionID)
		if err != nil {
			return err
		}

		return u.games.InsertGameQuestion(ctx, tx, specificGame.ID, specificQuestion.ID)
	})
}
