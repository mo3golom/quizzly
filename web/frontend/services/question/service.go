package question

import (
	"context"
	"errors"
	"github.com/a-h/templ"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/handlers"
	frontend "quizzly/web/frontend/templ"
	"quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
	"sort"
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
		return convertListToTempl(list, false, options.SelectIsEnabled, options.ActionsIsEnabled), nil
	default:
		return frontend.AdminPageComponent(
			listTitle,
			convertListToTempl(list, true, options.SelectIsEnabled, options.ActionsIsEnabled),
		), nil
	}
}

func convertListToTempl(in []model.Question, withHeader bool, withSelect bool, withActions bool) templ.Component {
	sort.Slice(in, func(i, j int) bool {
		return in[i].ID.String() < in[j].ID.String()
	})

	options := frontend_admin_question.Options{
		WithSelect:  withSelect,
		WithActions: withActions,
	}

	components := make([]templ.Component, 0, len(in)+1)
	if withHeader {
		components = append(components, frontendComponents.Header(
			listTitle,
			frontend_admin_question.ActionAddNewQuestion(),
		))
	}
	for _, item := range in {
		components = append(components, frontend_admin_question.QuestionListItem(
			handlers.Question{
				ID:      item.ID,
				ImageID: item.ImageID,
				Type:    item.Type,
				Text:    item.Text,
			},
			slices.SafeMap(item.AnswerOptions, func(ao model.AnswerOption) templ.Component {
				return frontend_admin_question.QuestionListItemAnswerOption(ao.Answer, ao.IsCorrect)
			}),
			options,
		))
	}

	return frontendComponents.Composition(components...)
}
