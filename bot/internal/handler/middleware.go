package handler

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

const RID = "X-Request-ID"

type middleware struct {
	lg                *slog.Logger
	rwMu              sync.RWMutex
	userSrv           userSrv
	messageSrv        messageSrv
	tgKeyboardSrv     tgKeyboardSrv
	requestLimitCount int
	activeRequest     map[int]struct{}
}

func newMiddleware(
	lg *slog.Logger,
	userSrv userSrv,
	messageSrv messageSrv,
	tgKeyboardSrv tgKeyboardSrv,
	requestLimitCount int,
) *middleware {
	return &middleware{
		lg:                lg,
		userSrv:           userSrv,
		messageSrv:        messageSrv,
		tgKeyboardSrv:     tgKeyboardSrv,
		requestLimitCount: requestLimitCount,
		activeRequest:     make(map[int]struct{}),
	}
}

func (m *middleware) setRIDAndLogDuration(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		const op = "handler.middleware.setRIDAndLogDuration"

		id := uuid.New().String()
		c.Set(RID, id)

		m.lg.Info(op, slog.String("RID", id))

		start := time.Now()
		err := next(c)
		duration := time.Since(start)

		m.lg.Info(
			op,
			slog.String("RID", id),
			slog.Any("seconds", duration.Seconds()),
		)

		return err
	}
}

func (m *middleware) setUserAndContext(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		const op = "handler.middleware.setUser"
		m.lg.Debug(op, slog.Any("RID", c.Get(RID)))

		ctx := context.Background()
		c.Set("ctx", ctx)

		user, err := m.userSrv.GetOrCreateByChatIDAndUsername(ctx, strconv.Itoa(int(c.Sender().ID)), c.Sender().Username)
		if err != nil {
			m.lg.Error(op, slog.String("error", err.Error()))
			return c.Send("Что-то пошло не так, напишите @eerzho")
		}
		c.Set("user", user)

		return next(c)
	}
}

func (m *middleware) spamLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		const op = "handler.middleware.spamLimit"
		m.lg.Debug(op, slog.Any("RID", c.Get(RID)))

		chatID := int(c.Sender().ID)
		isActive := m.isActiveRequest(chatID)

		if isActive {
			return c.Send("✨Пожалуйста, подождите✨", &telebot.SendOptions{ReplyTo: c.Message()})
		} else {
			m.setActiveRequest(chatID)
			defer m.delActiveRequest(chatID)
			return next(c)
		}
	}
}

func (m *middleware) setActiveRequest(chatID int) {
	m.rwMu.Lock()
	defer m.rwMu.Unlock()
	m.activeRequest[chatID] = struct{}{}
}

func (m *middleware) delActiveRequest(chatID int) {
	m.rwMu.Lock()
	defer m.rwMu.Unlock()
	delete(m.activeRequest, chatID)
}

func (m *middleware) isActiveRequest(chatID int) bool {
	m.rwMu.RLock()
	defer m.rwMu.RUnlock()
	_, ok := m.activeRequest[chatID]

	return ok
}

func (m *middleware) requestLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		const op = "handler.middleware.requestLimit"
		m.lg.Debug(op, slog.Any("RID", c.Get(RID)))
		ctx := c.Get("ctx").(context.Context)
		user := c.Get("user").(*model.User)

		if user.State != def.UserDefaultState {
			return next(c)
		}

		monthAgo := time.Now().AddDate(0, -1, 0)
		msgCount, err := m.messageSrv.CountByChatIDAndFromTime(ctx, user.ChatID, monthAgo)
		if err != nil {
			m.lg.Error(op, slog.String("error", err.Error()))
			return c.Send("✨Пожалуйста, повторите попытку позже✨")
		}

		if msgCount < m.requestLimitCount {
			return next(c)
		}

		if user.QuestionCount > 0 {
			err = next(c)
			if err != nil {
				return err
			}

			err = m.userSrv.DecreaseQuestionCount(ctx, user, 1)
			if err != nil {
				m.lg.Error(op, slog.String("error", err.Error()))
				return c.Send("✨Пожалуйста, повторите попытку позже✨")
			}

			return nil
		} else {
			return c.Send("✨Вы превысили лимит✨", &telebot.ReplyMarkup{
				InlineKeyboard: m.tgKeyboardSrv.OverLimit(ctx),
			})
		}
	}
}
