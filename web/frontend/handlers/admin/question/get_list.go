package question

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/auth"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/handlers"
	frontend "quizzly/web/frontend/templ"
	frontend_admin_question "quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
)

const (
	listTitle    = "Список вопросов"
	defaultLimit = 10
)

type (
	GetListData struct {
		PageNumber  *int64   `schema:"page_number"`
		IDs         []string `schema:"question_id"`
		WithSelect  bool     `schema:"with_select"`
		WithActions bool     `schema:"with_actions"`

		InContainer bool `schema:"in_container"`
	}

	GetHandler struct {
		uc contracts.QuestionUsecase
	}
)

func NewGetHandler(uc contracts.QuestionUsecase) *GetHandler {
	return &GetHandler{uc: uc}
}

func (h *GetHandler) Handle(_ http.ResponseWriter, request *http.Request, in GetListData) (templ.Component, error) {
	if !in.InContainer {
		return frontend.AdminPageComponent(
			listTitle,
			frontendComponents.Composition(
				frontendComponents.Header(
					listTitle,
					frontend_admin_question.ActionAddNewQuestion(),
				),
				frontend_admin_question.QuestionListContainer(frontend_admin_question.ContainerOptions{
					WithActions: true,
				}),
			),
		), nil
	}

	authContext := request.Context().(auth.Context)

	switch {
	case len(in.IDs) > 0:
		ids, err := slices.Map(in.IDs, func(id string) (uuid.UUID, error) {
			return uuid.Parse(id)
		})
		if err != nil {
			return nil, err
		}

		list, err := h.uc.GetByIDs(request.Context(), ids)
		if err != nil {
			return nil, err
		}

		return frontendComponents.Composition(
			convertListToTempl(list, in.WithSelect, in.WithActions),
		), nil
	default:
		page := int64(1)
		if in.PageNumber != nil {
			page = *in.PageNumber
		}

		result, err := h.uc.GetByAuthor(
			request.Context(),
			authContext.UserID(),
			page,
			defaultLimit,
		)
		if err != nil {
			return nil, err
		}

		return frontendComponents.Composition(
			convertListToTempl(result.Result, in.WithSelect, in.WithActions),
			frontendComponents.Pagination(
				page,
				result.TotalCount,
				defaultLimit,
			),
		), nil
	}
}

func convertListToTempl(in []model.Question, withSelect bool, withActions bool) templ.Component {
	templOptions := frontend_admin_question.Options{
		WithSelect:  withSelect,
		WithActions: withActions,
	}

	components := make([]templ.Component, 0, len(in)+1)
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
			templOptions,
		))
	}

	return frontendComponents.Composition(components...)
}
