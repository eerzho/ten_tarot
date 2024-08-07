package v1

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

const RID = "X-Request-ID"

type middleware struct {
	mu                sync.Mutex
	requestLimitCount int
	activeRequest     map[int]struct{}
	tgMessageService  tgMessageService
	tgKeyboardService tgKeyboardService
	tgUserService     tgUserService
}

func newMiddleware(
	tgMessageService tgMessageService,
	tgKeyboardService tgKeyboardService,
	tgUserService tgUserService,
) *middleware {
	return &middleware{
		requestLimitCount: 3,
		activeRequest:     make(map[int]struct{}),
		tgMessageService:  tgMessageService,
		tgKeyboardService: tgKeyboardService,
		tgUserService:     tgUserService,
	}
}

func (m *middleware) setRIDAndLogDuration(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "handler.telegram.v1.middleware.setRIDAndLogDuration"

		id := uuid.New().String()
		ctx.Set(RID, id)

		logger.Info(op, logger.Any("RID", id))

		start := time.Now()
		err := next(ctx)
		duration := time.Since(start)

		logger.Info(
			op,
			logger.Any("RID", id),
			logger.Any("SD", duration.Seconds()),
		)

		return err
	}
}

func (m *middleware) setUserAndContext(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "handler.telegram.v1.middleware.setUser"
		logger.Debug(op, logger.Any("id", ctx.Get(RID)))

		oc := context.Background()
		username := ctx.Sender().Username
		chatID := strconv.Itoa(int(ctx.Sender().ID))

		user, err := m.tgUserService.GetOrCreateByChatIDUsername(oc, chatID, username)
		if err != nil {
			logger.OPError(op, err)
			return ctx.Send("Что-то пошло не так, напишите @eerzho")
		}

		ctx.Set("oc", oc)
		ctx.Set("user", user)

		return next(ctx)
	}
}

func (m *middleware) spamLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "handler.telegram.v1.middleware.spamLimit"
		logger.Debug(op, logger.Any("id", ctx.Get(RID)))

		chatID := int(ctx.Sender().ID)
		isActive := m.isActiveRequest(chatID)

		if isActive {
			return ctx.Send(
				"✨Пожалуйста, подождите✨",
				&telebot.SendOptions{ReplyTo: ctx.Message()},
			)
		} else {
			m.setActiveRequest(chatID)
			defer func() {
				m.delActiveRequest(chatID)
			}()
			return next(ctx)
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
	return func(ctx telebot.Context) error {
		const op = "handler.telegram.v1.middleware.requestLimit"
		logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

		errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

		oc, ok := ctx.Get("oc").(context.Context)
		if !ok {
			logger.OPError(op, failure.ErrContextData)
			return ctx.Send(errTGMsg)
		}
		user, ok := ctx.Get("user").(*model.TGUser)
		if !ok {
			logger.OPError(op, failure.ErrContextData)
			return ctx.Send(errTGMsg)
		}

		if user.State != "" {
			return next(ctx)
		}

		monthAgo := time.Now().AddDate(0, -1, 0)
		msgCount, err := m.tgMessageService.CountByChatIDFromTime(oc, user.ChatID, monthAgo)
		if err != nil {
			logger.OPError(op, err)
			return ctx.Send(errTGMsg)
		}

		if msgCount < m.requestLimitCount {
			return next(ctx)
		}

		if user.QuestionCount > 0 {
			err = next(ctx)
			if err != nil {
				return err
			}
			err = m.tgUserService.DecreaseQC(oc, user, 1)
			if err != nil {
				logger.OPError(op, err)
				return ctx.Send(errTGMsg)
			}
			return nil
		} else {
			opt := telebot.ReplyMarkup{
				InlineKeyboard: m.tgKeyboardService.OverLimit(oc),
			}
			return ctx.Send("✨Вы превысили лимит✨", &opt)
		}
	}
}
