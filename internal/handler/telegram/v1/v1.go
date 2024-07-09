package v1

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/eerzho/ten_tarot/internal/service"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

const RID = "X-Request-ID"

func NewHandler(l logger.Logger, bot *telebot.Bot, tgUserService service.TGUser, tgMessageService service.TGMessage) {
	mv := newMiddleware(l, tgMessageService)
	bot.Use(mv.log)

	newCommand(l, bot, tgUserService)
	newMessage(l, mv, bot, tgMessageService, tgUserService)
}

type middleware struct {
	l                logger.Logger
	mu               sync.Mutex
	limit            int
	activeRequest    map[int64]struct{}
	tgMessageService service.TGMessage
}

func newMiddleware(l logger.Logger, tgMessageService service.TGMessage) *middleware {
	return &middleware{
		l:                l,
		limit:            10,
		activeRequest:    make(map[int64]struct{}),
		tgMessageService: tgMessageService,
	}
}

func (m *middleware) log(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		id := uuid.New().String()
		ctx.Set(RID, id)

		m.l.Info(fmt.Sprintf("start: %s", id))

		start := time.Now()
		err := next(ctx)
		duration := time.Since(start)

		m.l.Info(fmt.Sprintf("end: %s - %s", id, duration.String()))

		return err
	}
}

func (m *middleware) rateLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		userId := ctx.Message().Chat.ID
		exists := m.existsActiveRequest(userId)
		if exists {
			options := &telebot.SendOptions{ReplyTo: ctx.Message()}
			return ctx.Send("✨Пожалуйста, подождите✨", options)
		} else {
			m.setActiveRequest(userId)
			defer func() {
				m.delActiveRequest(userId)
			}()
			return next(ctx)
		}
	}
}

func (m *middleware) setActiveRequest(userId int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeRequest[userId] = struct{}{}
}

func (m *middleware) delActiveRequest(userId int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activeRequest, userId)
}

func (m *middleware) existsActiveRequest(userId int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.activeRequest[userId]

	return ok
}

func (m *middleware) dailyLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		st := time.Now().Add(-24 * time.Hour)

		count, err := m.tgMessageService.CountByTime(context.Background(), strconv.Itoa(int(ctx.Message().Chat.ID)), st)
		if err != nil {
			return next(ctx)
		}

		if count >= m.limit {
			opt := &telebot.SendOptions{ReplyTo: ctx.Message(), ParseMode: telebot.ModeMarkdown}
			return ctx.Send(fmt.Sprintf("✨Вы превысили лимит, попробуйте позже✨\n\n\n"+
				"🎁Скоро у вас появится возможность увеличить лимит, оплатив услугу или пригласив друзей🎁"), opt)
		}

		return next(ctx)
	}
}
