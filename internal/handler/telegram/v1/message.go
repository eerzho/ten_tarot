package v1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type (
	message struct {
		l                logger.Logger
		tgMessageService tgMessageService
		tgUserService    tgUserService
	}

	tgMessageService interface {
		CountByTime(ctx context.Context, chatID string, st time.Time) (int, error)
		Create(ctx context.Context, chatID, text string) (*model.TGMessage, error)
	}
)

func newMessage(l logger.Logger, mv *middleware, bot *telebot.Bot, tgMessageService tgMessageService, tgUserService tgUserService) *message {
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

	ctxB := context.Background()
	chatID := strconv.Itoa(int(ctx.Sender().ID))

	if _, err := m.tgUserService.Create(ctxB, chatID, ctx.Sender().Username); err != nil {
		m.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
	}

	opt := &telebot.SendOptions{ReplyTo: ctx.Message(), ParseMode: telebot.ModeMarkdown}

	waitMsg, err := ctx.Bot().Send(ctx.Sender(), "✨Пожалуйста, подождите✨", opt)
	if err != nil {
		m.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
	}

	tgMsg, err := m.tgMessageService.Create(ctxB, chatID, ctx.Message().Text)
	if err != nil {
		m.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
		if err = ctx.Send("✨Пожалуйста, повторите попытку позже✨", opt); err != nil {
			m.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
			return err
		}
		return err
	}

	if err = ctx.Bot().Delete(waitMsg); err != nil {
		m.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
	}

	if _, err = ctx.Bot().Send(ctx.Sender(), tgMsg.Answer, opt); err != nil {
		m.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
		return err
	}

	return nil
}
