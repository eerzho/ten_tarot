package v1

import (
	"sync"

	"github.com/eerzho/event_manager/pkg/logger"
	"gopkg.in/telebot.v3"
)

type middleware struct {
	l             logger.Logger
	mu            sync.Mutex
	activeRequest map[int64]struct{}
}

func newMiddleware(l logger.Logger) *middleware {
	return &middleware{
		l:             l,
		activeRequest: make(map[int64]struct{}),
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
			return ctx.Send("Вы отправляете сообщения слишком часто. Пожалуйста, подождите.", options)
		}

		defer func() {
			m.mu.Lock()
			delete(m.activeRequest, userId)
			m.mu.Unlock()
		}()

		return next(ctx)
	}
}
