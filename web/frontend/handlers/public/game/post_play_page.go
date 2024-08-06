package game

import (
	"context"
	"errors"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	"quizzly/internal/quizzly/contracts"
	"quizzly/internal/quizzly/model"
	"quizzly/web/frontend/handlers"
	frontend "quizzly/web/frontend/templ"
	frontendComponents "quizzly/web/frontend/templ/components"
	frontendPublicGame "quizzly/web/frontend/templ/public/game"
)

type (
	PostPlayPageData struct {
		GameID     uuid.UUID `schema:"id"`
		QuestionID uuid.UUID `schema:"question-id"`
		PlayerID   uuid.UUID `schema:"player-id"`
		Answers    []string  `schema:"answer"`
	}

	PostPlayPageHandler struct {
		gameUC    contracts.GameUsecase
		sessionUC contracts.SessionUsecase
		playerUC  contracts.PLayerUsecase
	}
)

func NewPostPlayPageHandler(
	gameUC contracts.GameUsecase,
	sessionUC contracts.SessionUsecase,
	playerUC contracts.PLayerUsecase,
) *PostPlayPageHandler {
	return &PostPlayPageHandler{
		gameUC:    gameUC,
		sessionUC: sessionUC,
		playerUC:  playerUC,
	}
}

func (h *PostPlayPageHandler) Handle(writer http.ResponseWriter, request *http.Request, in PostPlayPageData) (templ.Component, error) {
	game, err := h.gameUC.Get(request.Context(), in.GameID)
	if err != nil {
		return nil, err
	}
	if game == nil || game.Status == model.GameStatusFinished {
		return frontendPublicGame.StartPage(), nil
	}

	playerID := in.PlayerID
	if playerID == uuid.Nil {
		playerID = getPlayerID(request)
	}

	if playerID == uuid.Nil {
		return frontendPublicGame.StartPage(), nil
	}
	setPlayerID(writer, playerID)

	answerResult, err := h.sessionUC.AcceptAnswers(request.Context(), &contracts.AcceptAnswersIn{
		GameID:     game.ID,
		PlayerID:   playerID,
		QuestionID: in.QuestionID,
		Answers:    in.Answers,
	})
	if err != nil {
		return nil, err
	}

	session, err := h.sessionUC.GetCurrentState(request.Context(), in.GameID, playerID)
	if errors.Is(err, contracts.ErrQuestionQueueIsEmpty) {
		return h.finish(request.Context(), game, playerID, answerResult.IsCorrect)
	}
	if err != nil {
		return nil, err
	}

	specificQuestionColor := handlers.QuestionTypePublicColors
	answerOptions := make([]handlers.AnswerOption, 0, len(session.CurrentQuestion.AnswerOptions))
	for i, answerOption := range session.CurrentQuestion.AnswerOptions {
		answerOptions = append(answerOptions, handlers.AnswerOption{
			ID:    int64(answerOption.ID),
			Text:  answerOption.Answer,
			Color: frontend.ColorsMap[specificQuestionColor.AnswerOptionColors[i]][frontend.BgWithHoverColor],
		})
	}

	var playerName string
	players, err := h.playerUC.Get(request.Context(), []uuid.UUID{playerID})
	if err != nil {
		return nil, err
	}
	if len(players) > 0 {
		playerName = players[0].Name
	}

	var question templ.Component
	switch session.CurrentQuestion.Type {
	case model.QuestionTypeChoice, model.QuestionTypeOneOfChoice, model.QuestionTypeMultipleChoice:
		question = frontendPublicGame.Question(
			session.CurrentQuestion.ID,
			frontendPublicGame.QuestionBlock(session.CurrentQuestion.Text, session.CurrentQuestion.ImageID),
			frontendComponents.Composition(
				frontendPublicGame.AnswerChoiceDescription(session.CurrentQuestion.Type),
				frontendPublicGame.AnswerChoiceOptions(session.CurrentQuestion.Type, answerOptions),
			),
		)
	case model.QuestionTypeFillTheGap:
		question = frontendPublicGame.Question(
			session.CurrentQuestion.ID,
			frontendPublicGame.QuestionBlock(session.CurrentQuestion.Text, session.CurrentQuestion.ImageID),
			frontendPublicGame.AnswerTextInput(),
		)
	}

	return frontendPublicGame.QuestionForm(
		game.ID,
		playerID,
		frontendPublicGame.Header(game.Title),
		frontendComponents.GridLine(
			frontendPublicGame.Progress(&handlers.SessionProgress{
				Answered: int(session.Progress.Answered),
				Total:    int(session.Progress.Total),
			}),
			frontendPublicGame.Player(playerName),
		),
		question,
		frontendPublicGame.Answer(answerResult.IsCorrect),
	), nil
}

func (h *PostPlayPageHandler) finish(
	ctx context.Context,
	game *model.Game,
	playerID uuid.UUID,
	answerResult bool,
) (templ.Component, error) {
	err := h.sessionUC.Finish(context.Background(), game.ID, playerID)
	if err != nil {
		return nil, err
	}

	stats, err := h.sessionUC.GetStatistics(ctx, game.ID, playerID)
	if err != nil {
		return nil, err
	}

	return frontendComponents.Composition(
		frontendComponents.CompositionMD(
			frontendPublicGame.ResultHeader(game.Title),
			frontendPublicGame.ResultStatistics(
				&handlers.SessionStatistics{
					QuestionsCount:      int(stats.QuestionsCount),
					CorrectAnswersCount: int(stats.CorrectAnswersCount),
				},
			),
			frontendPublicGame.ActionRestartGame(game.ID),
		),
		frontendPublicGame.Answer(answerResult),
	), nil
}
