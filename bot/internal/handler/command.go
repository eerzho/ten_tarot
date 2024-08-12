package handler

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"log/slog"

	"gopkg.in/telebot.v3"
)

type command struct {
	lg           *slog.Logger
	userSrv      userSrv
	tgCommandSrv tgCommandSrv
}

func newCommand(
	bot *telebot.Bot,
	lg *slog.Logger,
	userSrv userSrv,
	tgCommandSrv tgCommandSrv,
) *command {
	c := &command{
		lg:           lg,
		userSrv:      userSrv,
		tgCommandSrv: tgCommandSrv,
	}

	bot.Handle(def.TGStartCommand, c.start)
	bot.Handle(def.TGDonateCommand, c.donate)
	bot.Handle(def.TGSupportCommand, c.support)
	bot.Handle(def.TGCancelCommand, c.cancel)

	c.setCommands(bot)

	return c
}

func (cmd *command) setCommands(bot *telebot.Bot) {
	const op = "handler.command.setCommands"
	cmd.lg.Debug(op)

	if err := bot.SetCommands(cmd.tgCommandSrv.GetCommands(context.Background())); err != nil {
		cmd.lg.Error(op, slog.String("error", err.Error()))
	}
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
	ctx := c.Get("ctx").(context.Context)
	user := c.Get("user").(*model.User)

	if err := cmd.userSrv.UpdateState(ctx, user, def.UserDonateState); err != nil {
		cmd.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send("Пожалуйста, введите сумму для доната в ⭐️")
}

func (cmd *command) support(c telebot.Context) error {
	const op = "handler.command.support"
	cmd.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := c.Get("ctx").(context.Context)
	user := c.Get("user").(*model.User)

	if err := cmd.userSrv.UpdateState(ctx, user, def.UserSupportState); err != nil {
		cmd.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send("Пожалуйста, напишите ваш запрос 🤕")
}

func (cmd *command) cancel(c telebot.Context) error {
	const op = "handler.command.cancel"
	cmd.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := c.Get("ctx").(context.Context)
	user := c.Get("user").(*model.User)

	if err := cmd.userSrv.UpdateState(ctx, user, def.UserDefaultState); err != nil {
		cmd.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send("Активная команда отменена, вы можете продолжить задавать вопросы боту 🤗")
}
