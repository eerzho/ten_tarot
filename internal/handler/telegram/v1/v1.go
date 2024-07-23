package v1

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

const RID = "X-Request-ID"

func NewHandler(bot *telebot.Bot, tgUserService tgUserService, tgMessageService tgMessageService) {
	mv := newMiddleware(tgMessageService)
	bot.Use(mv.log)

	newCommand(bot, tgUserService)
	newMessage(mv, bot, tgMessageService, tgUserService)
}

type middleware struct {
	mu               sync.Mutex
	limit            int
	activeRequest    map[int64]struct{}
	tgMessageService tgMessageService
}

func newMiddleware(tgMessageService tgMessageService) *middleware {
	return &middleware{
		limit:            10,
		activeRequest:    make(map[int64]struct{}),
		tgMessageService: tgMessageService,
	}
}

func (m *middleware) log(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		id := uuid.New().String()
		ctx.Set(RID, id)

		logger.Info(fmt.Sprintf("start: %s", id))

		start := time.Now()
		err := next(ctx)
		duration := time.Since(start)

		logger.Info(fmt.Sprintf("end: %s - %.4f sec.", id, duration.Seconds()))

		return err
	}
}

func (m *middleware) rateLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "./internal/handler/telegram/v1/middleware::rateLimit"

		userId := ctx.Message().Chat.ID
		exists := m.existsActiveRequest(userId)

		if exists {
			opt := &telebot.SendOptions{ReplyTo: ctx.Message()}
			if _, err := ctx.Bot().Send(ctx.Sender(), "✨Пожалуйста, подождите✨", opt); err != nil {
				logger.Error(fmt.Sprintf("%s - %s", op, err.Error()))
				return err
			}
			return nil
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
		const op = "./internal/handler/telegram/v1/middleware::dailyLimit"

		st := time.Now().Add(-24 * time.Hour)

		count, err := m.tgMessageService.CountByTime(context.Background(), strconv.Itoa(int(ctx.Message().Chat.ID)), st)
		if err != nil {
			logger.Error(fmt.Sprintf("%s - %s", op, err.Error()))
			return next(ctx)
		}

		if count >= m.limit {
			opt := &telebot.SendOptions{ReplyTo: ctx.Message(), ParseMode: telebot.ModeMarkdown}
			if _, err = ctx.Bot().Send(ctx.Sender(), fmt.Sprintf("✨Вы превысили лимит, попробуйте позже✨\n\n\n"+
				"🎁Скоро у вас появится возможность увеличить лимит, оплатив услугу или пригласив друзей🎁"), opt); err != nil {
				logger.Error(fmt.Sprintf("%s - %s", op, err.Error()))
				return err
			}
			return nil
		}

		return next(ctx)
	}
}
