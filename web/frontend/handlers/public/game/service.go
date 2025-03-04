package game

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/handlers"
	"quizzly/web/frontend/services/link"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

type (
	getCurrentStateIn struct {
		game       *model.Game
		player     *model.Player
		customName *string
	}

	service struct {
		sessionUC contracts.SessionUsecase

		linkService link.Service
	}
)

func (s *service) GetCurrentState(ctx context.Context, in *getCurrentStateIn, components ...templ.Component) (templ.Component, error) {
	if in.game.Status == model.GameStatusFinished {
		return frontendComponents.Redirect("/?warn=Игра уже завершена"), nil
	}
	if in.game.Status == model.GameStatusCreated {
		return frontendComponents.Redirect("/?warn=Игра еще не началась. Подождите немного или попросите автора запустить игру"), nil
	}

	if in.game.Settings.InputCustomName && !in.player.NameUserEntered && in.customName == nil {
		return frontendPublicGame.Page(
			frontendPublicGame.NamePage(in.game.Title, in.game.ID),
		), nil
	}

	session, err := s.sessionUC.GetCurrentState(ctx, in.game.ID, in.player.ID)
	if errors.Is(err, contracts.ErrQuestionQueueIsEmpty) {
		err = s.sessionUC.Finish(context.Background(), in.game.ID, in.player.ID)
		if err != nil {
			return nil, err
		}

		return frontendComponents.Composition(
			frontendPublicGame.ResultLinkInput(s.linkService.GameResultsLink(in.game.ID, in.player.ID)),
			frontendComponents.Composition(components...),
		), nil
	}
	if err != nil {
		return nil, err
	}

	if session.Status == model.SessionStatusFinished {
		return frontendComponents.Redirect(s.linkService.GameResultsLink(in.game.ID, in.player.ID)), nil
	}

	return frontendPublicGame.QuestionForm(
		in.game.ID,
		in.player.ID,
		frontendPublicGame.Header(in.game.Title),
		frontendComponents.GridLine(
			frontendPublicGame.Progress(&handlers.SessionProgress{
				Answered: int(session.Progress.Answered),
				Total:    int(session.Progress.Total),
			}),
			frontendPublicGame.Player(in.player.Name),
		),
		frontendPublicGame.Question(
			session.CurrentQuestion.ID,
			frontendPublicGame.QuestionBlock(session.CurrentQuestion.Text, session.CurrentQuestion.ImageID),
			frontendComponents.Composition(
				frontendPublicGame.AnswerChoiceDescription(session.CurrentQuestion.Type),
				getAnswerOptions(session.CurrentQuestion),
			),
		),
		frontendComponents.Composition(components...),
	), nil
}

func getAnswerOptions(question *model.Question) templ.Component {
	answerOptions := make([]handlers.AnswerOption, 0, len(question.AnswerOptions))
	for _, answerOption := range question.AnswerOptions {
		answerOptions = append(answerOptions, handlers.AnswerOption{
			ID:   int64(answerOption.ID),
			Text: answerOption.Answer,
		})
	}

	switch question.Type {
	case model.QuestionTypeFillTheGap:
		return frontendPublicGame.AnswerTextInput()
	case model.QuestionTypeMultipleChoice, model.QuestionTypeOneOfChoice:
		return frontendPublicGame.AnswerChoiceOptions(question.Type, answerOptions, true)
	default:
		return frontendPublicGame.AnswerChoiceOptions(question.Type, answerOptions)
	}
}

func gameTitle(game *model.Game) string {
	if game == nil {
		return "Игра не найдена"
	}

	title := fmt.Sprintf("Игра от %s", game.CreatedAt.Format("02.01.2006"))
	if game.Title != nil {
		title = fmt.Sprintf(`Игра "%s"`, *game.Title)
	}

	return title
}
