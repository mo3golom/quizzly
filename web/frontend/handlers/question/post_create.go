package question

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/structs"
	"quizzly/pkg/structs/collections/slices"
	"quizzly/web/frontend/services/question"
)

var (
	availableQuestionTypes = []string{
		string(model.QuestionTypeChoice),
		string(model.QuestionTypeMultipleChoice),
	}
)

type (
	NewPostData struct {
		QuestionText               string   `schema:"question_text"`
		QuestionType               string   `schema:"question_type"`
		QuestionMultipleChoiceType *string  `schema:"question_multiple_choice_type"`
		QuestionCorrectAnswer      []bool   `schema:"question_correct_answer"`
		QuestionAnswerOptionText   []string `schema:"question_answer_option_text"`
	}

	PostCreateHandler struct {
		uc      contracts.QuestionUsecase
		service question.Service
	}
)

func NewPostCreateHandler(uc contracts.QuestionUsecase, service question.Service) *PostCreateHandler {
	return &PostCreateHandler{uc: uc, service: service}
}

func (h *PostCreateHandler) Handle(_ http.ResponseWriter, request *http.Request, in NewPostData) (templ.Component, error) {
	converted, err := convert(&in)
	if err != nil {
		return nil, err
	}

	err = h.uc.Create(request.Context(), converted)
	if err != nil {
		return nil, err
	}

	return h.service.List(
		request.Context(),
		&question.Spec{
			AuthorID: structs.Pointer(uuid.New()),
		},
		&question.ListOptions{},
	)
}

func convert(in *NewPostData) (*model.Question, error) {
	clearIn(in)

	answerOptions := make([]model.AnswerOption, 0, len(in.QuestionAnswerOptionText))
	for i, text := range in.QuestionAnswerOptionText {
		answerOptions = append(answerOptions, model.AnswerOption{
			Answer:    text,
			IsCorrect: in.QuestionCorrectAnswer[i],
		})
	}

	if !slices.Contains(availableQuestionTypes, func(t string) bool {
		return t == in.QuestionType
	}) {
		return nil, errors.New("invalid question type")
	}

	questionType := model.QuestionType(in.QuestionType)
	if in.QuestionMultipleChoiceType != nil && *in.QuestionMultipleChoiceType == "one_of" {
		questionType = model.QuestionTypeOneOfChoice
	}

	return &model.Question{
		ID:            uuid.New(),
		Text:          in.QuestionText,
		Type:          questionType,
		AnswerOptions: answerOptions,
	}, nil
}

func clearIn(in *NewPostData) {
	if len(in.QuestionCorrectAnswer) <= 0 {
		return
	}

	questionCorrectAnswer := make([]bool, 0, len(in.QuestionCorrectAnswer))
	for i, value := range in.QuestionCorrectAnswer {
		if i == 0 || !value {
			questionCorrectAnswer = append(questionCorrectAnswer, value)
			continue
		}

		questionCorrectAnswer[i-1] = value
	}

	in.QuestionCorrectAnswer = questionCorrectAnswer
}
