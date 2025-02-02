package contracts

import (
	"context"
	"quizzly/internal/quizzly/model"

	"github.com/google/uuid"
)

type (
	CreateGameIn struct {
		AuthorID uuid.UUID
		Type     model.GameType
		Title    *string
		Settings model.GameSettings
	}

	GameUsecase interface {
		Create(ctx context.Context, in *CreateGameIn) (uuid.UUID, error)
		Update(ctx context.Context, in *model.Game) error
		Start(ctx context.Context, id uuid.UUID) error
		Finish(ctx context.Context, id uuid.UUID) error

		Get(ctx context.Context, id uuid.UUID) (*model.Game, error)
		GetByAuthor(ctx context.Context, authorID uuid.UUID) ([]model.Game, error)
		GetPublic(ctx context.Context) ([]model.Game, error)

		CreateQuestion(ctx context.Context, in *model.Question) error
		UpdateQuestion(ctx context.Context, in *model.Question) error
		DeleteQuestion(ctx context.Context, id uuid.UUID) error
		GetQuestions(ctx context.Context, gameID uuid.UUID) ([]model.Question, error)
	}
)
