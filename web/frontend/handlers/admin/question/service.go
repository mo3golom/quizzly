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

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

const (
	listTitle = "Список вопросов"
)

type service struct {
	uc contracts.GameUsecase
}

func (s *service) list(ctx context.Context, gameID uuid.UUID, inContainer bool, editable bool) (templ.Component, error) {
	if !inContainer {
		return frontend.AdminPageComponent(
			listTitle,
			frontendComponents.Composition(
				frontend_admin_question.QuestionListContainer(gameID, editable),
			),
		), nil
	}

	result, err := s.uc.GetQuestions(
		ctx,
		gameID,
	)
	if err != nil {
		return nil, err
	}

	return convertListToTempl(result, editable), nil
}

func convertListToTempl(in []model.Question, editable bool) templ.Component {
	components := make([]templ.Component, 0, len(in)+1)

	for i, question := range in {
		var actions []templ.Component
		if editable {
			actions = []templ.Component{
				frontend_admin_question.ActionDelete(question.ID),
			}
		}

		components = append(components, frontend_admin_question.QuestionListItem(
			i+1,
			handlers.Question{
				ID:      question.ID,
				ImageID: question.ImageID,
				Type:    question.Type,
				Text:    question.Text,
			},
			slices.SafeMap(question.AnswerOptions, func(ao model.AnswerOption) templ.Component {
				return frontend_admin_question.QuestionListItemAnswerOption(ao.Answer, ao.IsCorrect)
			}),
			actions,
		))

		if len(question.AnswerOptions) == 0 {
			continue
		}

	}

	if len(components) == 0 {
		return frontend_admin_question.NotFound()
	}

	return frontendComponents.Composition(components...)
}
