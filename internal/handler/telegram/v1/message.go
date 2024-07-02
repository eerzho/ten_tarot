package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/ten_tarot/internal/entity"
	"github.com/eerzho/ten_tarot/internal/service"
	"gopkg.in/telebot.v3"
)

type message struct {
	l                logger.Logger
	tgMessageService *service.TGMessage
	tgUserService    *service.TGUser
}

func newMessage(l logger.Logger, mv *middleware, bot *telebot.Bot, tgMessageService *service.TGMessage, tgUserService *service.TGUser) *message {
	m := &message{
		l:                l,
		tgMessageService: tgMessageService,
		tgUserService:    tgUserService,
	}

	bot.Handle(telebot.OnText, m.text, mv.rateLimit, mv.dailyLimit)

	return m
}

func (m *message) text(ctx telebot.Context) error {
	const op = "./internal/handler/telegram/v1/message::text"

	chatID := strconv.FormatInt(ctx.Sender().ID, 10)

	user := entity.TGUser{
		Username: ctx.Sender().Username,
		ChatID:   chatID,
	}
	if err := m.tgUserService.Create(context.Background(), &user); err != nil {
		m.l.Error(fmt.Errorf("%s: %w", op, err))
	}

	opt := &telebot.SendOptions{ReplyTo: ctx.Message(), ParseMode: telebot.ModeMarkdown}
	if err := ctx.Send("✨Пожалуйста, подождите✨", opt); err != nil {
		m.l.Error(fmt.Errorf("%s: %w", op, err))
	}

	msg := entity.TGMessage{
		Text:   ctx.Message().Text,
		ChatID: chatID,
	}
	if err := m.tgMessageService.Text(context.Background(), &msg); err != nil {
		m.l.Error(fmt.Errorf("%s: %w", op, err))
	}

	if msg.Answer != "" {
		if err := ctx.Send(msg.Answer, opt); err != nil {
			m.l.Error(fmt.Errorf("%s: %w", op, err))
			return err
		}
	} else {
		if err := ctx.Send("✨Пожалуйста, повторите попытку позже✨", opt); err != nil {
			m.l.Error(fmt.Errorf("%s: %w", op, err))
			return err
		}
	}

	return nil
}
