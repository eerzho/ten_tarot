package handler

import (
	"bot/internal/constant"
	"bot/internal/failure"
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
	mu                sync.Mutex
	requestLimitCount int
	activeRequest     map[int]struct{}
	tgMessageService  tgMessageService
	tgKeyboardService tgKeyboardService
	tgUserService     tgUserService
}

func newMiddleware(
	lg *slog.Logger,
	tgMessageService tgMessageService,
	tgKeyboardService tgKeyboardService,
	tgUserService tgUserService,
) *middleware {
	return &middleware{
		lg:                lg,
		requestLimitCount: 3,
		activeRequest:     make(map[int]struct{}),
		tgMessageService:  tgMessageService,
		tgKeyboardService: tgKeyboardService,
		tgUserService:     tgUserService,
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

func (m *middleware) setUser(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		const op = "handler.middleware.setUser"
		m.lg.Debug(op, slog.Any("RID", c.Get(RID)))
		ctx := context.Background()

		username := c.Sender().Username
		chatID := strconv.Itoa(int(c.Sender().ID))

		user, err := m.tgUserService.GetOrCreateByChatIDUsername(ctx, chatID, username)
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
			return c.Send(
				"✨Пожалуйста, подождите✨",
				&telebot.SendOptions{ReplyTo: c.Message()},
			)
		} else {
			m.setActiveRequest(chatID)
			defer func() {
				m.delActiveRequest(chatID)
			}()
			return next(c)
		}
	}
}

func (m *middleware) setActiveRequest(chatID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeRequest[chatID] = struct{}{}
}

func (m *middleware) delActiveRequest(chatID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activeRequest, chatID)
}

func (m *middleware) isActiveRequest(chatID int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.activeRequest[chatID]

	return ok
}

func (m *middleware) requestLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		const op = "handler.middleware.requestLimit"
		m.lg.Debug(op, slog.Any("RID", c.Get(RID)))
		ctx := context.Background()

		errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

		user, ok := c.Get("user").(*model.TGUser)
		if !ok {
			m.lg.Error(op, slog.String("error", failure.ErrContextData.Error()))
			return c.Send(errTGMsg)
		}

		if user.State != constant.UserDefaultState {
			return next(c)
		}

		monthAgo := time.Now().AddDate(0, -1, 0)
		msgCount, err := m.tgMessageService.CountByChatIDFromTime(ctx, user.ChatID, monthAgo)
		if err != nil {
			m.lg.Error(op, slog.String("error", err.Error()))
			return c.Send(errTGMsg)
		}

		if msgCount < m.requestLimitCount {
			return next(c)
		}

		if user.QuestionCount > 0 {
			err = next(c)
			if err != nil {
				return err
			}
			err = m.tgUserService.DecreaseQC(ctx, user, 1)
			if err != nil {
				m.lg.Error(op, slog.String("error", err.Error()))
				return c.Send(errTGMsg)
			}
			return nil
		} else {
			opt := telebot.ReplyMarkup{
				InlineKeyboard: m.tgKeyboardService.OverLimit(ctx),
			}
			return c.Send("✨Вы превысили лимит✨", &opt)
		}
	}
}
