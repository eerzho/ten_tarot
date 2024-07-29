package v1

import (
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type (
	command struct {
		tgUserService tgUserService
	}
)

func newCommand(bot *telebot.Bot, tgUserService tgUserService) *command {
	c := &command{
		tgUserService: tgUserService,
	}

	bot.Handle("/start", c.start)

	return c
}

func (c *command) start(ctx telebot.Context) error {
	const op = "handler.telegram.v1.command.start"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	if err := ctx.Send(
		"Я ваш личный Таролог и готов помочь вам получить ответы на любые вопросы. " +
			"Просто отправьте свой вопрос, и я сделаю расклад на Таро специально для вас.\n\n" +
			"✨Будьте готовы узнать, что приготовила для вас судьба✨",
	); err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}
