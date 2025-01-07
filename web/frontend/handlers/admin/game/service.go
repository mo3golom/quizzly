package game

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/services/link"
	frontendAdminGame "quizzly/web/frontend/templ/admin/game"
	frontendAdminQuestion "quizzly/web/frontend/templ/admin/question"
	frontendComponents "quizzly/web/frontend/templ/components"
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

	titleComponent := frontendAdminGame.Title(handlersGame)
	questionsComponent := frontendAdminQuestion.QuestionListContainer(game.ID)

	if game.Status == model.GameStatusCreated {
		titleComponent = frontendAdminGame.TitleInput(handlersGame)
		questionsComponent = frontendComponents.Composition(
			frontendAdminGame.ActionAddQuestion(),
			frontendAdminQuestion.QuestionListContainer(game.ID),
			frontendComponents.Modal(
				"addQuestionModal",
				"Добавить вопрос",
				frontendAdminQuestion.Form(
					game.ID,
					model.QuestionTypeChoice,
					frontendComponents.Composition(
						frontendAdminQuestion.QuestionImageInput(),
						frontendAdminQuestion.QuestionTextInput(),
					),
					frontendComponents.Composition(
						frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "orange", true),
						frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "pink", true),
						frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "amber", false),
						frontendAdminQuestion.AnswerChoiceInput(uuid.New(), "red", false),
					),
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
			frontendComponents.Composition(settingsComponents...),
		),
		frontendAdminGame.Invite(s.linkService.GameLink(game.ID, request)),
		frontendComponents.Tabs(
			uuid.New(),
			frontendComponents.Tab{
				Name:    "Вопросы",
				Content: questionsComponent,
			},
			frontendComponents.Tab{
				Name:    "Участники",
				Content: frontendAdminGame.SessionListContainer(game.ID),
			},
		),
	), nil
}
