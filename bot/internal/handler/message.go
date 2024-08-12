package handler

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"errors"
	"log/slog"

	"gopkg.in/telebot.v3"
)

type (
	message struct {
		lg                *slog.Logger
		userSrv           userSrv
		messageSrv        messageSrv
		tgInvoiceSrv      tgInvoiceSrv
		supportRequestSrv supportRequestSrv
	}
)

func newMessage(
	bot *telebot.Bot,
	lg *slog.Logger,
	mdw *middleware,
	userSrv userSrv,
	messageSrv messageSrv,
	tgInvoiceSrv tgInvoiceSrv,
	supportRequestSrv supportRequestSrv,
) *message {
	m := &message{
		lg:                lg,
		userSrv:           userSrv,
		messageSrv:        messageSrv,
		tgInvoiceSrv:      tgInvoiceSrv,
		supportRequestSrv: supportRequestSrv,
	}

	bot.Handle(telebot.OnText, m.text, mdw.spamLimit, mdw.requestLimit)

	return m
}

func (m *message) text(c telebot.Context) error {
	const op = "handler.message.text"
	m.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := c.Get("ctx").(context.Context)
	user := c.Get("user").(*model.User)

	switch user.State {
	case def.UserDonateState:
		return m.generateInvoice(c, ctx, user)
	case def.UserSupportState:
		return m.saveRequest(c, ctx, user)
	default:
		return m.generateAnswer(c, ctx, user)
	}
}

func (m *message) generateInvoice(c telebot.Context, ctx context.Context, user *model.User) error {
	const op = "handler.message.generateInvoice"
	m.lg.Debug(
		op,
		slog.Any("user", user),
	)

	tgInvoice, err := m.tgInvoiceSrv.CreateDonateInvoice(ctx, user, c.Message().Text)
	if err != nil {
		if errors.Is(err, def.ErrInvalidType) {
			m.lg.Warn(op, slog.String("error", err.Error()))
			return c.Send("Пожалуйста, введите только цифру 0️⃣-9️⃣")
		}
		m.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send(tgInvoice)
}

func (m *message) saveRequest(c telebot.Context, ctx context.Context, user *model.User) error {
	const op = "handler.message.saveRequest"
	m.lg.Debug(
		op,
		slog.Any("user", user),
	)

	if _, err := m.supportRequestSrv.CreateByUserQuestion(ctx, user, c.Message().Text); err != nil {
		m.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send("Спасибо за ваш запрос 😁")
}

func (m *message) generateAnswer(c telebot.Context, ctx context.Context, user *model.User) error {
	const op = "handler.message.generateAnswer"
	m.lg.Debug(
		op,
		slog.Any("user", user),
	)

	opt := telebot.SendOptions{
		ReplyTo: c.Message(),
	}

	waitTGMsg, err := c.Bot().Send(c.Sender(), "✨Пожалуйста, подождите✨", &opt)
	if err != nil {
		m.lg.Warn(op, slog.String("error", err.Error()))
	}

	msg, err := m.messageSrv.CreateByChatIDAndUserQuestion(ctx, user.ChatID, c.Message().Text)
	if err != nil {
		m.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	if err = c.Bot().Delete(waitTGMsg); err != nil {
		m.lg.Warn(op, slog.String("error", err.Error()))
	}

	return c.Send(msg.BotAnswer, &opt)
}
