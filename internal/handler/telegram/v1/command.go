package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/service"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type command struct {
	l             logger.Logger
	tgUserService *service.TGUser
}

func newCommand(l logger.Logger, bot *telebot.Bot, tgUserService *service.TGUser) *command {
	c := &command{
		l:             l,
		tgUserService: tgUserService,
	}

	bot.Handle("/start", c.start)

	return c
}

func (c *command) start(ctx telebot.Context) error {
	const op = "./internal/handler/telegram/v1/command::start"

	_, err := c.tgUserService.Create(context.Background(), strconv.Itoa(int(ctx.Sender().ID)), ctx.Sender().Username)
	if err != nil {
		c.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
	}

	if err = ctx.Send("Я ваш личный Таролог и готов помочь вам получить ответы на любые вопросы. " +
		"Просто отправьте свой вопрос, и я сделаю расклад на Таро специально для вас.\n\n" +
		"✨Будьте готовы узнать, что приготовила для вас судьба✨"); err != nil {
		c.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
		return err
	}

	return nil
}
