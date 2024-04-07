package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"quizzly/pkg/logger"
)

const (
	webhookRoute = "/updates"
)

type DefaultGateway struct {
	*handler

	client externalClient
	token  string
}

func NewGateway(
	c externalClient,
	log logger.Logger,
	token string,
) *DefaultGateway {
	gateway := &DefaultGateway{
		client: c,
		token:  token,
		handler: &handler{
			client:          c,
			errorLogHandler: &ErrorLogHandler{log: log},
		},
	}

	return gateway
}

func (g *DefaultGateway) Run(ctx context.Context, config Config) {
	tgConfig := tgbotapi.NewUpdate(config.Offset)
	tgConfig.Timeout = config.Timeout
	tgConfig.Limit = config.Limit

	var updates tgbotapi.UpdatesChannel
	if config.Webhook != nil && config.Webhook.Enable {
		url := fmt.Sprintf("%s/%s", config.Webhook.Host, g.token)
		if !config.Webhook.Debug {
			webhook, err := tgbotapi.NewWebhook(fmt.Sprintf("%s%s", url, webhookRoute))
			if err != nil {
				panic(err)
			}

			_, err = g.client.Request(webhook)
			if err != nil {
				panic(err)
			}

			info, err := g.client.GetWebhookInfo()
			if err != nil {
				panic(err)
			}

			if info.LastErrorDate != 0 {
				panic(fmt.Sprintf("failed to set webhook: %s", info.LastErrorMessage))
			}
		}

		updates = g.client.ListenForWebhook(webhookRoute)
		go func(addr string) {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				panic(err)
			}
		}(url)

		fmt.Printf("bot started, webhook: %s", url)
	} else {
		updates = g.client.GetUpdatesChan(tgConfig)
	}

	g.handler.handle(ctx, updates)
}
