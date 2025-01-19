package question

import (
	"context"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/handlers"
	frontend "quizzly/web/frontend/templ"
	frontend_admin_question "quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

const (
	listTitle = "Список вопросов"
)

var (
	v2ReleaseDate, _ = time.Parse("15:04 02.01.2006", "00:01 01.01.2025")
)

type service struct {
	uc contracts.GameUsecase
}

func (s *service) list(ctx context.Context, gameID uuid.UUID, inContainer bool) (templ.Component, error) {
	if !inContainer {
		return frontend.AdminPageComponent(
			listTitle,
			frontendComponents.Composition(
				frontend_admin_question.QuestionListContainer(gameID),
			),
		), nil
	}

	specificGame, err := s.uc.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	result, err := s.uc.GetQuestions(
		ctx,
		gameID,
	)
	if err != nil {
		return nil, err
	}

	if specificGame.CreatedAt.Before(v2ReleaseDate) {
		return convertListToTemplOld(result), nil
	}

	return convertListToTemplNew(result), nil
}

func convertListToTemplNew(in *model.QuestionMap) templ.Component {
	components := make([]templ.Component, 0, in.Len()+1)
	question := in.GetFirst()
	index := 0

	for question != nil {
		index++
		components = append(components, frontend_admin_question.QuestionListItem(
			index,
			handlers.Question{
				ID:      question.ID,
				ImageID: question.ImageID,
				Type:    question.Type,
				Text:    question.Text,
			},
			slices.SafeMap(question.AnswerOptions, func(ao model.AnswerOption) templ.Component {
				return frontend_admin_question.QuestionListItemAnswerOption(ao.Answer, ao.IsCorrect)
			}),
		))

		if len(question.AnswerOptions) == 0 {
			question = nil
			continue
		}

		question, _ = in.GetNextQuestion(question.ID, question.AnswerOptions[0].ID)
	}

	if len(components) == 0 {
		return frontend_admin_question.NotFound()
	}

	return frontendComponents.Composition(components...)
}

func convertListToTemplOld(in *model.QuestionMap) templ.Component {
	components := make([]templ.Component, 0, in.Len()+1)
	for i, question := range in.Values() {
		components = append(components, frontend_admin_question.QuestionListItem(
			i,
			handlers.Question{
				ID:      question.ID,
				ImageID: question.ImageID,
				Type:    question.Type,
				Text:    question.Text,
			},
			slices.SafeMap(question.AnswerOptions, func(ao model.AnswerOption) templ.Component {
				return frontend_admin_question.QuestionListItemAnswerOption(ao.Answer, ao.IsCorrect)
			}),
		))
	}

	if len(components) == 0 {
		return frontend_admin_question.NotFound()
	}

	return frontendComponents.Composition(components...)
}
