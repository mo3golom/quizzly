package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"quizzly/pkg/logger"
	"quizzly/pkg/variables"
)

type Configuration struct {
	Gateway Gateway
}

func NewConfiguration(log logger.Logger, token string) *Configuration {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	if variables.AppEnvironment() != variables.EnvironmentProd {
		bot.Debug = true
	}

	return &Configuration{
		Gateway: NewGateway(bot, log, token),
	}
}
