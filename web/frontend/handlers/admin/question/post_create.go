package question

import (
	"bytes"
	"errors"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/pkg/files"
	"quizzly/pkg/structs/collections/slices"
	"strings"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

const (
	questionImageName = "question_image"
)

var (
	availableQuestionTypes = []string{
		string(model.QuestionTypeChoice),
		string(model.QuestionTypeMultipleChoice),
		string(model.QuestionTypeOneOfChoice),
		string(model.QuestionTypeFillTheGap),
	}
)

type (
	NewPostData struct {
		QuestionText               string    `schema:"question_text"`
		QuestionType               string    `schema:"question_type"`
		QuestionMultipleChoiceType *string   `schema:"question_multiple_choice_type"`
		QuestionCorrectAnswer      []bool    `schema:"question_correct_answer"`
		QuestionAnswerOptionText   []string  `schema:"question_answer_option_text"`
		GameID                     uuid.UUID `schema:"game_id"`
	}

	PostCreateHandler struct {
		uc      contracts.GameUsecase
		images  files.Manager
		service *service
	}
)

func NewPostCreateHandler(
	uc contracts.GameUsecase,
	images files.Manager,
) *PostCreateHandler {
	return &PostCreateHandler{
		uc:      uc,
		images:  images,
		service: &service{uc: uc},
	}
}

func (h *PostCreateHandler) Handle(_ http.ResponseWriter, request *http.Request, in NewPostData) (templ.Component, error) {
	converted, err := convert(&in)
	if err != nil {
		return nil, err
	}

	image, err := findQuestionImage(request)
	if err != nil {
		return nil, err
	}
	if image != nil {
		image.Name = strings.ReplaceAll(strings.TrimSpace(image.Name), " ", "_")
		err = h.images.Upload(request.Context(), image)
		if err != nil {
			return nil, err
		}
		converted.ImageID = &image.Name
	}

	err = h.uc.CreateQuestion(request.Context(), converted)
	if err != nil {
		return nil, err
	}

	return h.service.list(request.Context(), in.GameID, true, true)
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
		GameID:        in.GameID,
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

		questionCorrectAnswer[len(questionCorrectAnswer)-1] = value
	}

	in.QuestionCorrectAnswer = questionCorrectAnswer
}

func findQuestionImage(r *http.Request) (*files.UploadFile, error) {
	file, handler, err := r.FormFile(questionImageName)
	if errors.Is(err, http.ErrMissingFile) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return &files.UploadFile{
		Data: &buffer,
		Name: handler.Filename,
		Size: handler.Size,
	}, nil
}
