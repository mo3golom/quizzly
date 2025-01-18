package game

import (
	"context"
	"github.com/google/uuid"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/internal/quizzly/repositories/game"
	"quizzly/internal/quizzly/repositories/session"
	"quizzly/pkg/structs"
	"quizzly/pkg/transactional"
)

type Usecase struct {
	games    game.Repository
	sessions session.Repository
	template transactional.Template
}

func NewUsecase(
	games game.Repository,
	sessions session.Repository,
	template transactional.Template,
) contracts.GameUsecase {
	return &Usecase{
		games:    games,
		sessions: sessions,
		template: template,
	}
}

func (u *Usecase) Create(ctx context.Context, in *contracts.CreateGameIn) (uuid.UUID, error) {
	id := uuid.New()

	return id, u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.games.Upsert(
			ctx,
			tx,
			&model.Game{
				ID:       id,
				AuthorID: in.AuthorID,
				Status:   model.GameStatusCreated,
				Type:     model.GameTypeAsync,
				Title:    in.Title,
				Settings: in.Settings,
			},
		)
	})

}

func (u *Usecase) Update(ctx context.Context, in *model.Game) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.games.Upsert(ctx, tx, in)
	})
}

func (u *Usecase) Start(ctx context.Context, id uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		specificGames, err := u.games.GetBySpecWithTx(ctx, tx, &game.Spec{
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
		return u.games.Upsert(ctx, tx, &specificGame)
	})
}

func (u *Usecase) Finish(ctx context.Context, id uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		specificGames, err := u.games.GetBySpecWithTx(ctx, tx, &game.Spec{
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
		return u.games.Upsert(ctx, tx, &specificGame)
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

	err := u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.games.InsertQuestion(ctx, tx, in)
	})
	if err != nil {
		return err
	}

	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		questions, err := u.games.GetQuestionsBySpecWithTx(ctx, tx, &game.QuestionsSpec{
			GameID: &in.GameID,
			Order: &game.Order{
				Field:     "created_at",
				Direction: "desc",
			},
		})
		if err != nil {
			return err
		}

		if len(questions) <= 1 {
			return nil
		}

		previousQuestion := &questions[1] // get previous question
		for i := range previousQuestion.AnswerOptions {
			if previousQuestion.AnswerOptions[i].NextQuestionID != nil {
				continue
			}

			previousQuestion.AnswerOptions[i].NextQuestionID = &in.ID
		}

		return u.games.UpdateQuestion(ctx, tx, previousQuestion)
	})
}

func (u *Usecase) UpdateQuestion(ctx context.Context, in *model.Question) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.games.UpdateQuestion(ctx, tx, in)
	})
}

func (u *Usecase) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return u.template.Execute(ctx, func(tx transactional.Tx) error {
		return u.games.DeleteQuestion(ctx, tx, id)
	})
}

func (u *Usecase) GetQuestions(ctx context.Context, gameID uuid.UUID) (*model.QuestionMap, error) {
	result, err := u.games.GetQuestionsBySpec(ctx, &game.QuestionsSpec{
		GameID: &gameID,
		Order: &game.Order{
			Field:     "created_at",
			Direction: "asc",
		},
	})
	if err != nil {
		return nil, err
	}

	return model.NewQuestionMap(result), nil
}
