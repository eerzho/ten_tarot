package v1

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eerzho/event_manager/internal/entity"
	"github.com/eerzho/event_manager/internal/service"
	"github.com/eerzho/event_manager/pkg/logger"
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

	bot.Handle(telebot.OnText, m.text, mv.rateLimit)

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

	msg := entity.TGMessage{
		Text:   ctx.Message().Text,
		ChatID: chatID,
	}
	if err := m.tgMessageService.Text(context.Background(), &msg); err != nil {
		m.l.Error(fmt.Errorf("%s: %w", op, err))
	}
	defer func() {
		if msg.File != "" {
			if err := os.Remove(msg.File); err != nil {
				m.l.Error(fmt.Errorf("%s: %w", op, err))
			}
		}
	}()

	// send google calendar link
	if msg.Answer != "" {
		if err := ctx.Send(msg.Answer, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown, ReplyTo: ctx.Message()}); err != nil {
			m.l.Error(fmt.Errorf("%s: %w", op, err))
			return err
		}
	}

	// send file for apple calendar
	if msg.File != "" {
		file := telebot.FromDisk(msg.File)
		doc := telebot.Document{File: file, FileName: strings.Replace(strings.ToLower(ctx.Message().Text), " ", "_", -1) + ".ics", MIME: "text/calendar"}
		if err := ctx.Send(&doc); err != nil {
			m.l.Error(fmt.Errorf("%s: %w", op, err))
			return err
		}
	}

	return nil
}
