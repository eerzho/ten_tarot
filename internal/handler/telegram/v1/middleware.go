package v1

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/ten_tarot/internal/service"
	"gopkg.in/telebot.v3"
)

type middleware struct {
	l                logger.Logger
	mu               sync.Mutex
	activeRequest    map[int64]struct{}
	tgMessageService *service.TGMessage
}

func newMiddleware(l logger.Logger, tgMessageService *service.TGMessage) *middleware {
	return &middleware{
		l:                l,
		activeRequest:    make(map[int64]struct{}),
		tgMessageService: tgMessageService,
	}
}

func (m *middleware) rateLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		userId := ctx.Message().Chat.ID

		m.mu.Lock()
		_, exists := m.activeRequest[userId]
		if !exists {
			m.activeRequest[userId] = struct{}{}
		}
		m.mu.Unlock()

		if exists {
			options := &telebot.SendOptions{ReplyTo: ctx.Message()}
			return ctx.Send("✨Пожалуйста, подождите✨", options)
		}

		defer func() {
			m.mu.Lock()
			delete(m.activeRequest, userId)
			m.mu.Unlock()
		}()

		return next(ctx)
	}
}

func (m *middleware) dailyLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	const op = "./internal/handler/telegram/v1/middleware::dailyLimit"
	return func(ctx telebot.Context) error {
		userID := strconv.Itoa(int(ctx.Message().Chat.ID))
		count, err := m.tgMessageService.CountByDay(context.Background(), userID)
		m.l.Debug("count", count)
		if err != nil {
			m.l.Error(fmt.Errorf("%s: %w", op, err))
			return next(ctx)
		}

		if count >= 10 {
			resetTime := time.Until(time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour))
			hours := int(resetTime.Hours())
			minutes := int(resetTime.Minutes()) % 60
			options := &telebot.SendOptions{ReplyTo: ctx.Message(), ParseMode: telebot.ModeMarkdown}
			return ctx.Send(fmt.Sprintf("Вы превысили лимит, вы сможете отправить сообщение через `%dh %dm`", hours, minutes), options)
		}

		return next(ctx)
	}
}
