package v1

import (
	"context"
	"strconv"

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

	ctxB := context.Background()
	chatID := strconv.Itoa(int(ctx.Sender().ID))

	if _, err := m.tgUserService.GetOrCreateByChatIDUsername(
		ctxB,
		chatID, ctx.Sender().Username,
	); err != nil {
		logger.OPError(op, err)
	}

	opt := &telebot.SendOptions{
		ReplyTo:   ctx.Message(),
		ParseMode: telebot.ModeMarkdown,
	}

	waitMsg, err := ctx.Bot().Send(
		ctx.Sender(),
		"✨Пожалуйста, подождите✨",
		opt,
	)
	if err != nil {
		logger.OPWarn(op, err)
	}

	tgMsg, err := m.tgMessageService.CreateByChatIDUQ(
		ctxB,
		chatID,
		ctx.Message().Text,
	)
	if err != nil {
		logger.OPError(op, err)
		if err = ctx.Send("✨Пожалуйста, повторите попытку позже✨", opt); err != nil {
			logger.OPError(op, err)
			return err
		}
		return err
	}

	if err = ctx.Bot().Delete(waitMsg); err != nil {
		logger.OPWarn(op, err)
	}

	if err = ctx.Send(
		tgMsg.BotAnswer,
		opt,
	); err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}
