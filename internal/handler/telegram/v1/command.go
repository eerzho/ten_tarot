package v1

import (
	"context"
	"strconv"

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

	_, err := c.tgUserService.GetOrCreateByChatIDUsername(
		context.Background(),
		strconv.Itoa(int(ctx.Sender().ID)),
		ctx.Sender().Username,
	)
	if err != nil {
		logger.OPError(op, err)
	}

	if _, err = ctx.Bot().Send(
		ctx.Sender(),
		"Я ваш личный Таролог и готов помочь вам получить ответы на любые вопросы. "+
			"Просто отправьте свой вопрос, и я сделаю расклад на Таро специально для вас.\n\n"+
			"✨Будьте готовы узнать, что приготовила для вас судьба✨"); err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}
