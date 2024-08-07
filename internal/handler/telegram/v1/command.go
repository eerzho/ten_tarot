package v1

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/constant"
	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
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
	bot.Handle("/donate", c.donate)
	bot.Handle("/support", c.support)
	bot.Handle("/cancel", c.cancel)

	commands := []telebot.Command{
		{Text: "/donate", Description: "Пожертвовать в развитие проекта"},
		{Text: "/support", Description: "Связаться с разработчиком"},
		{Text: "/cancel", Description: "Отменить активную команду"},
	}

	if err := bot.SetCommands(commands); err != nil {
		logger.Fatal(err.Error())
	}

	return c
}

func (c *command) start(ctx telebot.Context) error {
	const op = "handler.telegram.v1.command.start"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	return ctx.Send(
		"Я ваш личный Таролог и готов помочь вам получить ответы на любые вопросы. " +
			"Просто отправьте свой вопрос, и я сделаю расклад на Таро специально для вас.\n\n" +
			"✨Будьте готовы узнать, что приготовила для вас судьба✨",
	)
}

func (c *command) donate(ctx telebot.Context) error {
	const op = "handler.telegram.v1.command.donate"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := ctx.Get("user").(*model.TGUser)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}
	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}

	if err := c.tgUserService.UpdateState(oc, user, constant.DonateState); err != nil {
		logger.OPError(op, err)
		return ctx.Send(errTGMsg)
	}

	return ctx.Send("Пожалуйста, введите сумму для доната в ⭐️")
}

func (c *command) support(ctx telebot.Context) error {
	const op = "handler.telegram.v1.command.support"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := ctx.Get("user").(*model.TGUser)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}
	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}

	if err := c.tgUserService.UpdateState(oc, user, constant.SupportState); err != nil {
		logger.OPError(op, err)
		return ctx.Send(errTGMsg)
	}

	return ctx.Send("Пожалуйста, напишите ваш запрос 🤕")
}

func (c *command) cancel(ctx telebot.Context) error {
	const op = "handler.telegram.v1.command.cancel"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := ctx.Get("user").(*model.TGUser)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}
	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}

	if err := c.tgUserService.UpdateState(oc, user, ""); err != nil {
		logger.OPError(op, err)
		return ctx.Send(errTGMsg)
	}

	return ctx.Send("Активная команда отменена, вы можете продолжить задавать вопросы боту 🤗")
}
