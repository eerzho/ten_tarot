package v1

import (
	"context"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type (
	message struct {
		tgMessageService tgMessageService
		tgUserService    tgUserService
	}
)

func newMessage(mv *middleware, bot *telebot.Bot, tgMessageService tgMessageService, tgUserService tgUserService) *message {
	m := &message{
		tgMessageService: tgMessageService,
		tgUserService:    tgUserService,
	}

	bot.Handle(telebot.OnText, m.text, mv.spamLimit, mv.requestLimit)

	return m
}

func (m *message) text(ctx telebot.Context) error {
	const op = "handler.telegram.v1.message.text"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	chatID := strconv.Itoa(int(ctx.Sender().ID))
	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}

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
		return ctx.Send(errTGMsg)
	}

	if err = ctx.Bot().Delete(waitMsg); err != nil {
		logger.OPWarn(op, err)
	}

	return ctx.Send(tgMsg.BotAnswer, &opt)
}
