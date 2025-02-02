package game

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs"
)

type Usecase struct {
	games    game.Repository
	sessions session.Repository
	trm      trm.Manager
}

func NewUsecase(
	games game.Repository,
	sessions session.Repository,
	trm trm.Manager,
) contracts.GameUsecase {
	return &Usecase{
		games:    games,
		sessions: sessions,
		trm:      trm,
	}
}

func (u *Usecase) Create(ctx context.Context, in *contracts.CreateGameIn) (uuid.UUID, error) {
	id := uuid.New()

	return id, u.games.Upsert(
		ctx,
		&model.Game{
			ID:       id,
			AuthorID: in.AuthorID,
			Status:   model.GameStatusCreated,
			Type:     model.GameTypeAsync,
			Title:    in.Title,
			Settings: in.Settings,
		},
	)

}

func (u *Usecase) Update(ctx context.Context, in *model.Game) error {
	return u.games.Upsert(ctx, in)
}

func (u *Usecase) Start(ctx context.Context, id uuid.UUID) error {
	return u.trm.Do(ctx, func(ctx context.Context) error {
		specificGames, err := u.games.GetBySpec(ctx, &game.Spec{
			IDs: []uuid.UUID{id},
		})
		if err != nil {
			return err
		}
		if len(specificGames) == 0 {
			return contracts.ErrGameNotFound
		}

		specificGame := specificGames[0]

		questions, err := u.games.GetQuestionsBySpec(ctx, &game.QuestionsSpec{GameID: &specificGame.ID})
		if err != nil {
			return err
		}
		if len(questions) == 0 {
			return contracts.ErrEmptyQuestions
		}

		specificGame.Status = model.GameStatusStarted
		return u.games.Upsert(ctx, &specificGame)
	})
}

func (u *Usecase) Finish(ctx context.Context, id uuid.UUID) error {
	return u.trm.Do(ctx, func(ctx context.Context) error {
		specificGames, err := u.games.GetBySpec(ctx, &game.Spec{
			IDs: []uuid.UUID{id},
		})
		if err != nil {
			return err
		}
		if len(specificGames) == 0 {
			return contracts.ErrGameNotFound
		}

		specificGame := specificGames[0]
		specificGame.Status = model.GameStatusFinished
		return u.games.Upsert(ctx, &specificGame)
	})
}

func (u *Usecase) Get(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	specificGames, err := u.games.GetBySpec(ctx, &game.Spec{
		IDs: []uuid.UUID{id},
	})
	if err != nil {
		return nil, err
	}
	if len(specificGames) == 0 {
		return nil, contracts.ErrGameNotFound
	}

	return &specificGames[0], nil
}

func (u *Usecase) GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]model.Game, error) {
	return u.games.GetBySpec(ctx, &game.Spec{
		AuthorID: &authorID,
	})
}

func (u *Usecase) GetPublic(ctx context.Context) ([]model.Game, error) {
	return u.games.GetBySpec(ctx, &game.Spec{
		IsPrivate: structs.Pointer(false),
		Statuses:  []model.GameStatus{model.GameStatusStarted},
		Limit:     10,
	})
}

func (u *Usecase) CreateQuestion(ctx context.Context, in *model.Question) error {
	if len(in.AnswerOptions) == 0 {
		return contracts.ErrEmptyAnswerOptions
	}

	if in.ID == uuid.Nil {
		in.ID = uuid.New()
	}

	return u.games.InsertQuestion(ctx, in)
}

func (u *Usecase) UpdateQuestion(ctx context.Context, in *model.Question) error {
	return u.games.UpdateQuestion(ctx, in)
}

func (u *Usecase) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return u.games.DeleteQuestion(ctx, id)
}

func (u *Usecase) GetQuestions(ctx context.Context, gameID uuid.UUID) ([]model.Question, error) {
	result, err := u.games.GetQuestionsBySpec(ctx, &game.QuestionsSpec{
		GameID: &gameID,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
