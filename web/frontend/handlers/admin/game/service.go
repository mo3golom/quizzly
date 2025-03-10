package game

import (
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/services/link"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
	frontendAdminQuestion "quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

const (
	listUrl = "/admin/game/list"
)

var (
	settings = []setting{
		{
			slug: "is_private",
			text: "Частная игра",
			hint: "если вкл., то игра будет доступна только по ссылке. Игра не будет отображаться в списке \"новых игр\" на главной и странице результатов",
			value: func(settings *model.GameSettings) bool {
				return settings.IsPrivate
			},
		},
		{
			slug: "shuffle_questions",
			text: "Перемешать вопросы",
			hint: "каждый игрок будет видеть вопросы в случайном порядке. Если игрок не ответил на вопрос, он увидит его снова при следующем входе в игру.",
			value: func(settings *model.GameSettings) bool {
				return settings.ShuffleQuestions
			},
		},
		{
			slug: "shuffle_answers",
			text: "Перемешать ответы в вопросе",
			hint: "ответы на каждый вопрос будут перемешиваться для каждого игрока. Это помогает избежать запоминания игроками правильного порядка ответов.",
			value: func(settings *model.GameSettings) bool {
				return settings.ShuffleAnswers
			},
		},
		{
			slug: "show_right_answers",
			text: "Показывать правильный ответ в случае неудачи",
			hint: "при неправильном ответе игрока на экран результатов будет выводиться правильный ответ. Обратите внимание, что кнопка \"играть снова\" всегда активна, так что игрок может запомнить правильные ответы и пройти викторину без ошибок во второй раз.",
			value: func(settings *model.GameSettings) bool {
				return settings.ShowRightAnswers
			},
		},
		{
			slug: "input_custom_name",
			text: "Игрок должен ввести имя перед игрой",
			hint: "перед началом игры игрок должен ввести имя/псевдоним. При этом экран не будет показан при повторной игре",
			value: func(settings *model.GameSettings) bool {
				return settings.InputCustomName
			},
		},
	}
)

type service struct {
	uc          contracts.GameUsecase
	linkService link.Service
}

func (s *service) getGamePage(request *http.Request, gameID uuid.UUID) (templ.Component, error) {
	game, err := s.uc.Get(request.Context(), gameID)
	if err != nil {
		return nil, err
	}
	handlersGame := convertModelGameToHandlersGame(game)

	titleComponent := frontendAdminGame.Title(game.Title)
	questionsComponent := frontendAdminQuestion.QuestionListContainer(game.ID, game.Status == model.GameStatusCreated)

	if game.Status == model.GameStatusCreated {
		titleComponent = frontendAdminGame.TitleInput(game.ID, game.Title)
		questionsComponent = frontendComponents.Composition(
			frontendComponents.CompositionMB4(frontendAdminGame.ActionAddQuestion()),
			questionsComponent,
			frontendComponents.Modal(
				"addQuestionModal",
				"Добавить вопрос",
				frontendComponents.Tabs(
					uuid.New(),
					frontendComponents.Tab{
						Name:    "Один ответ",
						Content: singleChoiceQuestionForm(game.ID),
					},
					frontendComponents.Tab{
						Name:    "Несколько ответов",
						Content: multipleChoiceQuestionForm(game.ID),
					},
					frontendComponents.Tab{
						Name:    "Ввод слова",
						Content: fillTheGapQuestionForm(game.ID),
					},
				),
			),
		)
	}

	settingsComponents := make([]templ.Component, 0, len(settings))
	for _, item := range settings {
		if game.Status != model.GameStatusCreated && !item.value(&game.Settings) {
			continue
		}
		if game.Status != model.GameStatusCreated {
			settingsComponents = append(settingsComponents, frontendAdminGame.SettingBadge(item.text, item.hint))
			continue
		}

		settingsComponents = append(settingsComponents, frontendAdminGame.SettingBadge(
			item.text,
			item.hint,
			frontendAdminGame.SettingToggle(
				game.ID,
				item.slug,
				item.value(&game.Settings),
			),
		))
	}

	return frontendAdminGame.Page(
		frontendComponents.BackLink(listUrl),
		frontendAdminGame.Header(
			handlersGame,
			titleComponent,
		),
		frontendAdminGame.Invite(s.linkService.GameLink(game.ID, request)),
		frontendComponents.Tabs(
			uuid.New(),
			frontendComponents.Tab{
				Name: "Игра",
				Content: frontendComponents.Composition(
					frontendComponents.CompositionMB4(settingsComponents...),
					questionsComponent,
				),
			},
			frontendComponents.Tab{
				Name:    "Участники",
				Content: frontendAdminGame.SessionListContainer(game.ID),
			},
		),
	), nil
}

func singleChoiceQuestionForm(gameID uuid.UUID) templ.Component {
	return frontendAdminQuestion.Form(
		gameID,
		model.QuestionTypeChoice,
		frontendComponents.Composition(
			frontendAdminQuestion.QuestionImageInput(),
			frontendAdminQuestion.QuestionTextInput(),
		),
		frontendComponents.Composition(
			frontendAdminQuestion.AnswerChoiceInput(0, uuid.New(), true),
			frontendAdminQuestion.AnswerChoiceInput(1, uuid.New(), true),
			frontendAdminQuestion.AnswerChoiceInput(2, uuid.New(), false),
			frontendAdminQuestion.AnswerChoiceInput(3, uuid.New(), false),
		),
	)
}

func multipleChoiceQuestionForm(gameID uuid.UUID) templ.Component {
	return frontendAdminQuestion.Form(
		gameID,
		model.QuestionTypeMultipleChoice,
		frontendComponents.Composition(
			frontendAdminQuestion.QuestionImageInput(),
			frontendAdminQuestion.QuestionTextInput(),
			frontendAdminQuestion.QuestionMultipleChoiceOption(),
		),
		frontendComponents.Composition(
			frontendAdminQuestion.AnswerChoiceInput(0, uuid.New(), true, true),
			frontendAdminQuestion.AnswerChoiceInput(1, uuid.New(), true, true),
			frontendAdminQuestion.AnswerChoiceInput(2, uuid.New(), false, true),
			frontendAdminQuestion.AnswerChoiceInput(3, uuid.New(), false, true),
		),
	)
}

func fillTheGapQuestionForm(gameID uuid.UUID) templ.Component {
	return frontendAdminQuestion.Form(
		gameID,
		model.QuestionTypeFillTheGap,
		frontendComponents.Composition(
			frontendAdminQuestion.QuestionImageInput(),
			frontendAdminQuestion.QuestionTextInput(),
		),
		frontendComponents.Composition(
			frontendAdminQuestion.AnswerTextInput(),
		),
	)
}
