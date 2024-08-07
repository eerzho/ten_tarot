package v1

import (
	"context"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/constant"
	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type (
	message struct {
		tgMessageService      tgMessageService
		tgUserService         tgUserService
		tgInvoiceService      tgInvoiceService
		supportRequestService supportRequestService
	}
)

func newMessage(
	mv *middleware,
	bot *telebot.Bot,
	tgMessageService tgMessageService,
	tgUserService tgUserService,
	tgInvoiceService tgInvoiceService,
	supportRequestService supportRequestService,
) *message {
	m := &message{
		tgMessageService:      tgMessageService,
		tgUserService:         tgUserService,
		tgInvoiceService:      tgInvoiceService,
		supportRequestService: supportRequestService,
	}

	bot.Handle(telebot.OnText, m.text, mv.spamLimit, mv.requestLimit)

	return m
}

func (m *message) text(ctx telebot.Context) error {
	const op = "handler.telegram.v1.message.text"
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

	switch user.State {
	case constant.DonateState:
		return m.generateInvoice(ctx, oc, user)
	case constant.SupportState:
		return m.saveRequest(ctx, oc, user)
	default:
		return m.generateAnswer(ctx, oc)
	}
}

func (m *message) generateInvoice(ctx telebot.Context, oc context.Context, user *model.TGUser) error {
	const op = "handler.telegram.v1.message.generateInvoice"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	starsCount, err := strconv.Atoi(ctx.Message().Text)
	if err != nil {
		logger.OPWarn(op, err)
		return ctx.Send("Пожалуйста, введите только цифру 0️⃣-9️⃣")
	}

	chatID := strconv.Itoa(int(ctx.Sender().ID))

	tgInvoice, err := m.tgInvoiceService.CreateByChatIDSC(oc, chatID, starsCount)
	if err != nil {
		logger.OPError(op, err)
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

	if err = m.tgUserService.UpdateState(oc, user, ""); err != nil {
		logger.OPError(op, err)
		return ctx.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return ctx.Send(&invoice)
}

func (m *message) saveRequest(ctx telebot.Context, oc context.Context, user *model.TGUser) error {
	const op = "handler.telegram.v1.message.saveRequest"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	if _, err := m.supportRequestService.CreateByUserQuestion(oc, user, ctx.Message().Text); err != nil {
		logger.OPError(op, err)
		return ctx.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return ctx.Send("Спасибо за ваш запрос 😁")
}

func (m *message) generateAnswer(ctx telebot.Context, oc context.Context) error {
	const op = "handler.telegram.v1.message.generateAnswer"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	chatID := strconv.Itoa(int(ctx.Sender().ID))
	opt := telebot.SendOptions{
		ReplyTo: ctx.Message(),
	}

	waitMsg, err := ctx.Bot().Send(ctx.Sender(), "✨Пожалуйста, подождите✨", &opt)
	if err != nil {
		logger.OPWarn(op, err)
	}

	tgMsg, err := m.tgMessageService.CreateByChatIDUQ(oc, chatID, ctx.Message().Text)
	if err != nil {
		logger.OPError(op, err)
		return ctx.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	if err = ctx.Bot().Delete(waitMsg); err != nil {
		logger.OPWarn(op, err)
	}

	return ctx.Send(tgMsg.BotAnswer, &opt)
}
