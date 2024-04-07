package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"quizzly/pkg/helper"
	"quizzly/pkg/telegram/model"
)

type (
	handler struct {
		client externalClient

		commandHandlers map[string]Handler
		errorHandlers   []ErrorHandler
		messageHandlers []Handler
		middlewares     []Middleware

		errorLogHandler ErrorHandler
	}
)

func (h *handler) RegisterMiddleware(middleware ...Middleware) {
	copyMiddleware := helper.CopySlice[Middleware](h.middlewares)
	copyMiddleware = append(copyMiddleware, middleware...)

	h.middlewares = copyMiddleware
}

func (h *handler) RegisterCommandHandler(handler ...CommandHandler) {
	copyCommandHandlers := helper.CopyMap[string, Handler](h.commandHandlers)

	for _, h := range handler {
		copyCommandHandlers[h.Command()] = h
		for _, alias := range h.Aliases() {
			copyCommandHandlers[alias] = h
		}
	}
	h.commandHandlers = copyCommandHandlers
}

func (h *handler) RegisterErrorHandler(handler ...ErrorHandler) {
	copyErrorHandlers := helper.CopySlice[ErrorHandler](h.errorHandlers)
	copyErrorHandlers = append(copyErrorHandlers, handler...)

	h.errorHandlers = copyErrorHandlers
}

func (h *handler) RegisterHandler(handler ...Handler) {
	copyMessageHandlers := helper.CopySlice[Handler](h.messageHandlers)
	copyMessageHandlers = append(copyMessageHandlers, handler...)

	h.messageHandlers = copyMessageHandlers
}

func (h *handler) handle(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	senderImpl := &sender{
		client: h.client,
	}

	for update := range updates {
		go func(in tgbotapi.Update) {
			request := model.NewRequest(in)
			senderAdapterImpl := &senderAdapter{
				sender: senderImpl,
				chatID: request.Chat.ID,
			}
			err := h.handleUpdate(ctx, &request, senderAdapterImpl)
			if err == nil {
				return
			}

			h.handleError(ctx, err, &request, senderAdapterImpl)
		}(update)
	}
}

func (h *handler) handleUpdate(ctx context.Context, request *model.Request, sender Sender) error {
	// MIDDLEWARE
	for _, middleware := range h.middlewares {
		err := middleware.Handle(ctx, request)
		if err != nil {
			return err
		}
	}

	if request.CallbackID != nil {
		err := sender.SendCallback(*request.CallbackID)
		if err != nil {
			return err
		}
	}

	// COMMAND
	command, ok := h.determineCommand(request)
	if ok {
		return command.Handle(ctx, request, sender)
	}

	// MESSAGE HANDLER
	for _, handler := range h.messageHandlers {
		err := handler.Handle(ctx, request, sender)
		if err == nil {
			continue
		}

		return err
	}

	if len(h.messageHandlers) == 0 {
		_, err := sender.Send(
			model.
				NewResponse().
				SetText("Unknown command"),
		)
		return err
	}

	return nil
}

func (h *handler) handleError(ctx context.Context, err error, request *model.Request, sender Sender) {
	var ok bool
	for _, errHandler := range h.errorHandlers {
		if errHandler == nil {
			continue
		}

		ok = errHandler.Handle(ctx, err, request, sender)
		if !ok {
			continue
		}

		return
	}

	if !ok {
		h.errorLogHandler.Handle(ctx, err, request, sender)
	}
}

func (h *handler) determineCommand(request *model.Request) (Handler, bool) {
	command, ok := h.commandHandlers[request.Command]
	return command, ok
}
