package handler

import (
	"bot/internal/constant"
	"bot/internal/failure"
	"bot/internal/model"
	"context"
	"gopkg.in/telebot.v3"
	"log/slog"
	"strconv"
)

type (
	message struct {
		lg                    *slog.Logger
		tgMessageService      tgMessageService
		tgUserService         tgUserService
		tgInvoiceService      tgInvoiceService
		supportRequestService supportRequestService
	}
)

func newMessage(
	bot *telebot.Bot,
	lg *slog.Logger,
	mdw *middleware,
	tgMessageService tgMessageService,
	tgUserService tgUserService,
	tgInvoiceService tgInvoiceService,
	supportRequestService supportRequestService,
) *message {
	m := &message{
		lg:                    lg,
		tgMessageService:      tgMessageService,
		tgUserService:         tgUserService,
		tgInvoiceService:      tgInvoiceService,
		supportRequestService: supportRequestService,
	}

	bot.Handle(telebot.OnText, m.text, mdw.spamLimit, mdw.requestLimit)

	return m
}

func (m *message) text(c telebot.Context) error {
	const op = "handler.message.text"
	m.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := context.Background()

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := c.Get("user").(*model.TGUser)
	if !ok {
		m.lg.Error(op, slog.String("error", failure.ErrContextData.Error()))
		return c.Send(errTGMsg)
	}

	switch user.State {
	case constant.UserDonateState:
		return m.generateInvoice(c, ctx, user)
	case constant.UserSupportState:
		return m.saveRequest(c, ctx, user)
	default:
		return m.generateAnswer(c, ctx)
	}
}

func (m *message) generateInvoice(c telebot.Context, ctx context.Context, user *model.TGUser) error {
	const op = "handler.message.generateInvoice"
	m.lg.Debug(op, slog.Any("RID", c.Get(RID)))

	starsCount, err := strconv.Atoi(c.Message().Text)
	if err != nil {
		m.lg.Warn(op, slog.String("error", err.Error()))
		return c.Send("Пожалуйста, введите только цифру 0️⃣-9️⃣")
	}

	chatID := strconv.Itoa(int(c.Sender().ID))

	tgInvoice, err := m.tgInvoiceService.CreateByChatIDSC(ctx, chatID, starsCount)
	if err != nil {
		m.lg.Error(op, slog.String("error", err.Error()))
		return err
	}

	invoice := telebot.Invoice{
		Title:       "Благодарим за поддержку!",
		Description: "Ваш вклад поможет развивать проект и продвигать его дальше!",
		Payload:     tgInvoice.ID,
		Currency:    "XTR",
		Prices: []telebot.Price{
			{
				Label:  "Поддержка проекта",
				Amount: tgInvoice.StarsCount,
			},
		},
	}

	if err = m.tgUserService.UpdateState(ctx, user, constant.UserDefaultState); err != nil {
		m.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send(&invoice)
}

func (m *message) saveRequest(c telebot.Context, ctx context.Context, user *model.TGUser) error {
	const op = "handler.message.saveRequest"
	m.lg.Debug(op, slog.Any("RID", c.Get(RID)))

	if _, err := m.supportRequestService.CreateByUserQuestion(ctx, user, c.Message().Text); err != nil {
		m.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send("Спасибо за ваш запрос 😁")
}

func (m *message) generateAnswer(c telebot.Context, ctx context.Context) error {
	const op = "handler.message.generateAnswer"
	m.lg.Debug(op, slog.Any("RID", c.Get(RID)))

	chatID := strconv.Itoa(int(c.Sender().ID))
	opt := telebot.SendOptions{
		ReplyTo: c.Message(),
	}

	waitMsg, err := c.Bot().Send(c.Sender(), "✨Пожалуйста, подождите✨", &opt)
	if err != nil {
		m.lg.Warn(op, slog.String("error", err.Error()))
	}

	tgMsg, err := m.tgMessageService.CreateByChatIDUQ(ctx, chatID, c.Message().Text)
	if err != nil {
		m.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	if err = c.Bot().Delete(waitMsg); err != nil {
		m.lg.Warn(op, slog.String("error", err.Error()))
	}

	return c.Send(tgMsg.BotAnswer, &opt)
}
