package handler

import (
	"bot/internal/constant"
	"bot/internal/failure"
	"bot/internal/model"
	"context"
	"gopkg.in/telebot.v3"
	"log/slog"
)

type (
	command struct {
		lg            *slog.Logger
		tgUserService tgUserService
	}
)

func newCommand(bot *telebot.Bot, lg *slog.Logger, tgUserService tgUserService) *command {
	c := &command{
		lg:            lg,
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
		lg.Error("handler.command", slog.String("error", err.Error()))
	}

	return c
}

func (cmd *command) start(c telebot.Context) error {
	const op = "handler.command.start"
	cmd.lg.Debug(op, slog.Any("RID", c.Get(RID)))

	return c.Send(
		"Я ваш личный Таролог и готов помочь вам получить ответы на любые вопросы. " +
			"Просто отправьте свой вопрос, и я сделаю расклад на Таро специально для вас.\n\n" +
			"✨Будьте готовы узнать, что приготовила для вас судьба✨",
	)
}

func (cmd *command) donate(c telebot.Context) error {
	const op = "handler.command.donate"
	cmd.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := context.Background()

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := c.Get("user").(*model.TGUser)
	if !ok {
		cmd.lg.Error(op, slog.String("error", failure.ErrContextData.Error()))
		return c.Send(errTGMsg)
	}

	if err := cmd.tgUserService.UpdateState(ctx, user, constant.UserDonateState); err != nil {
		cmd.lg.Error(op, slog.String("error", err.Error()))
		return c.Send(errTGMsg)
	}

	return c.Send("Пожалуйста, введите сумму для доната в ⭐️")
}

func (cmd *command) support(c telebot.Context) error {
	const op = "handler.command.support"
	cmd.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := context.Background()

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := c.Get("user").(*model.TGUser)
	if !ok {
		cmd.lg.Error(op, slog.String("error", failure.ErrContextData.Error()))
		return c.Send(errTGMsg)
	}

	if err := cmd.tgUserService.UpdateState(ctx, user, constant.UserSupportState); err != nil {
		cmd.lg.Error(op, slog.String("error", err.Error()))
		return c.Send(errTGMsg)
	}

	return c.Send("Пожалуйста, напишите ваш запрос 🤕")
}

func (cmd *command) cancel(c telebot.Context) error {
	const op = "handler.command.cancel"
	cmd.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := context.Background()

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := c.Get("user").(*model.TGUser)
	if !ok {
		cmd.lg.Error(op, slog.String("error", failure.ErrContextData.Error()))
		return c.Send(errTGMsg)
	}

	if err := cmd.tgUserService.UpdateState(ctx, user, constant.UserDefaultState); err != nil {
		cmd.lg.Error(op, slog.String("error", err.Error()))
		return c.Send(errTGMsg)
	}

	return c.Send("Активная команда отменена, вы можете продолжить задавать вопросы боту 🤗")
}
