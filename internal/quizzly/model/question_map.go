package model

import (
	"errors"
	"math/rand/v2"

	"github.com/google/uuid"
)

var (
	ErrQuestionNotFound = errors.New("question not found")
)

type (
	QuestionMap struct {
		first uuid.UUID
		items map[uuid.UUID]QuestionMapItem
	}

	QuestionMapItem struct {
		question Question
		next     map[AnswerOptionID]uuid.UUID
	}
)

func NewQuestionMap(in []Question) *QuestionMap {
	result := &QuestionMap{}

	allNextMap := make(map[uuid.UUID]bool, len(in))
	items := make(map[uuid.UUID]QuestionMapItem, len(in))

	for _, item := range in {
		item := item

		nextMap := make(map[AnswerOptionID]uuid.UUID, len(item.AnswerOptions))
		for _, answerOption := range item.AnswerOptions {
			if answerOption.NextQuestionID == nil {
				continue
			}

			nextMap[answerOption.ID] = *answerOption.NextQuestionID
			allNextMap[*answerOption.NextQuestionID] = true
		}

		items[item.ID] = QuestionMapItem{
			question: item,
			next:     nextMap,
		}
	}
	result.items = items

	for key := range items {
		if _, ok := allNextMap[key]; ok {
			continue
		}

		result.first = key
		break
	}

	return result
}

func (q *QuestionMap) GetFirst() *Question {
	question, _ := q.GetQuestion(q.first)
	return question
}

func (q *QuestionMap) GetQuestion(questionID uuid.UUID) (*Question, error) {
	item, ok := q.items[questionID]
	if !ok {
		return nil, ErrQuestionNotFound
	}

	return &item.question, nil
}

func (q *QuestionMap) GetRandomQuestion() *Question {
	if q.Empty() {
		return nil
	}

	keys := q.GetIDs()
	question, _ := q.GetQuestion(keys[rand.IntN(len(keys))])
	return question
}

func (q *QuestionMap) GetNextQuestion(questionID uuid.UUID, answerOptionID AnswerOptionID) (*Question, error) {
	item, ok := q.items[questionID]
	if !ok {
		return nil, ErrQuestionNotFound
	}

	nextQuestionID, ok := item.next[answerOptionID]
	if !ok {
		return nil, ErrQuestionNotFound
	}

	return q.GetQuestion(nextQuestionID)
}

func (q *QuestionMap) GetIDs() []uuid.UUID {
	if q.Empty() {
		return nil
	}

	keys := make([]uuid.UUID, 0, len(q.items))
	for k := range q.items {
		keys = append(keys, k)
	}

	return keys
}

func (q *QuestionMap) Empty() bool {
	return len(q.items) == 0
}

func (q *QuestionMap) Len() int {
	return len(q.items)
}

func (q *QuestionMap) Keys() []uuid.UUID {
	result := make([]uuid.UUID, 0, len(q.items))
	for key := range q.items {
		result = append(result, key)
	}

	return result
}

func (q *QuestionMap) Values() []Question {
	result := make([]Question, 0, len(q.items))
	for _, item := range q.items {
		result = append(result, item.question)
	}

	return result
}
