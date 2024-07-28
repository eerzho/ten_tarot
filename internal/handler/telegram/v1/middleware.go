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

type middleware struct {
	mu                sync.Mutex
	requestLimitCount int
	activeRequest     map[int64]struct{}
	tgMessageService  tgMessageService
	tgButtonService   tgButtonService
	tgUserService     tgUserService
}

func newMiddleware(
	tgMessageService tgMessageService,
	tgButtonService tgButtonService,
	tgUserService tgUserService,
) *middleware {
	return &middleware{
		requestLimitCount: 3,
		activeRequest:     make(map[int64]struct{}),
		tgMessageService:  tgMessageService,
		tgButtonService:   tgButtonService,
		tgUserService:     tgUserService,
	}
}

func (m *middleware) setRIDAndLogDuration(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "handler.http.v1.middleware.setRIDAndLogDuration"

		id := uuid.New().String()
		ctx.Set(RID, id)

		logger.Info(op, logger.Any("id", id))

		start := time.Now()
		err := next(ctx)
		duration := time.Since(start)

		logger.Info(
			op,
			logger.Any("id", id),
			logger.Any("duration in seconds", duration.Seconds()),
		)

		return err
	}
}

func (m *middleware) spamLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "handler.http.v1.middleware.requestLimit"

		userId := ctx.Message().Chat.ID
		isActive := m.isActiveRequest(userId)

		if isActive {
			if err := ctx.Send(
				"✨Пожалуйста, подождите✨",
				&telebot.SendOptions{ReplyTo: ctx.Message()},
			); err != nil {
				logger.OPError(op, err)
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

func (m *middleware) isActiveRequest(userId int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.activeRequest[userId]

	return ok
}

func (m *middleware) requestLimit(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "handler.http.v1.middleware.requestLimit"

		c := context.Background()
		monthAgo := time.Now().AddDate(0, -1, 0)

		chatID := strconv.Itoa(int(ctx.Sender().ID))
		count, err := m.tgMessageService.CountByChatIDFromTime(c, chatID, monthAgo)
		if err != nil {
			logger.OPError(op, err)
			return next(ctx)
		}

		if count >= m.requestLimitCount {
			user, err := m.tgUserService.GetOrCreateByChatIDUsername(c, chatID, ctx.Sender().Username)
			if err != nil {
				logger.OPError(op, err)
				return err
			}

			if user.QuestionCount == 0 {
				opt := telebot.ReplyMarkup{
					InlineKeyboard: m.tgButtonService.OverLimit(c),
				}

				if _, err := ctx.Bot().Send(ctx.Sender(), "✨Вы превысили лимит✨", &opt); err != nil {
					logger.Error(fmt.Sprintf("%s - %s", op, err.Error()))
					return err
				}

				return nil
			}

			user.QuestionCount--
			if _, err := m.tgUserService.UpdateByChatIDQC(c, user.ChatID, user.QuestionCount); err != nil {
				logger.OPError(op, err)
				return err
			}

			return next(ctx)
		}

		return next(ctx)
	}
}
