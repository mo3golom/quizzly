package question

import (
	"context"
	"errors"
	"github.com/a-h/templ"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	frontend "quizzly/web/frontend/templ"
	"quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
)

const (
	listTitle = "Список вопросов"
)

type DefaultService struct {
	uc contracts.QuestionUsecase
}

func NewService(uc contracts.QuestionUsecase) *DefaultService {
	return &DefaultService{
		uc: uc,
	}
}

func (s *DefaultService) List(ctx context.Context, spec *Spec, options *ListOptions) (templ.Component, error) {
	var list []model.Question
	var err error
	switch {
	case len(spec.QuestionIDs) > 0:
		list, err = s.uc.GetByIDs(ctx, spec.QuestionIDs)
		if err != nil {
			return nil, err
		}
	case spec.AuthorID != nil:
		list, err = s.uc.GetByAuthor(ctx, *spec.AuthorID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("empty spec")
	}

	switch options.Type {
	case ListTypeCompact:
		return convertListToTempl(list, false, options.SelectIsEnabled), nil
	default:
		return frontend.AdminPageComponent(
			listTitle,
			convertListToTempl(list, true, options.SelectIsEnabled),
		), nil
	}
}

func convertListToTempl(in []model.Question, withHeader bool, withSelect bool) templ.Component {
	options := frontend_question.Options{
		WithSelect: withSelect,
	}

	components := make([]templ.Component, 0, len(in)+1)
	if withHeader {
		components = append(components, frontendComponents.Header(
			listTitle,
			frontend_question.ActionAddNewQuestion(),
		))
	}
	for _, item := range in {
		components = append(components, frontend_question.QuestionListItem(
			item.ID,
			item.Text,
			slices.SafeMap(item.AnswerOptions, func(ao model.AnswerOption) templ.Component {
				return frontend_question.QuestionListItemAnswerOption(ao.Answer, ao.IsCorrect)
			}),
			options,
		))
	}

	return frontendComponents.Composition(components...)
}
