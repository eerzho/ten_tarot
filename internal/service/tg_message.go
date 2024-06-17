package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/eerzho/event_manager/internal/entity"
	"github.com/eerzho/event_manager/internal/failure"
	"github.com/eerzho/event_manager/pkg/logger"
)

type (
	TGMessageRepo interface {
		All(ctx context.Context, chatID string, page, count int) ([]entity.TGMessage, error)
		Create(ctx context.Context, message *entity.TGMessage) error
	}

	TGMessage struct {
		l                     logger.Logger
		repo                  TGMessageRepo
		tgUserService         *TGUser
		eventService          *Event
		googleCalendarService *GoogleCalendar
		appleCalendar         *AppleCalendar
	}
)

func NewTGMessage(
	l logger.Logger,
	repo TGMessageRepo,
	tgUserService *TGUser,
	eventService *Event,
	googleCalendarService *GoogleCalendar,
	appleCalendar *AppleCalendar,
) *TGMessage {
	return &TGMessage{
		l:                     l,
		repo:                  repo,
		tgUserService:         tgUserService,
		eventService:          eventService,
		googleCalendarService: googleCalendarService,
		appleCalendar:         appleCalendar,
	}
}

func (t *TGMessage) All(ctx context.Context, chatID string, page, count int) ([]entity.TGMessage, error) {
	const op = "./internal/service.tg_message::All"

	messages, err := t.repo.All(ctx, chatID, page, count)
	if err != nil {
		t.l.Error(fmt.Errorf("%s: %w", op, err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}

func (t *TGMessage) Text(ctx context.Context, message *entity.TGMessage) error {
	const op = "./internal/service/tg_message::Text"

	defer func() {
		if message.Answer != "" {
			if err := t.repo.Create(ctx, message); err != nil {
				t.l.Error(fmt.Errorf("%s: %w", op, err))
			}
		}
	}()

	var event entity.Event
	if err := t.eventService.CreateFromText(ctx, &event, message.Text); err != nil {
		if errors.Is(err, failure.ErrValidation) && event.Message != "" {
			message.Answer = event.Message
			return nil
		}
		t.l.Error(fmt.Errorf("%s: %w", op, err))
		return fmt.Errorf("%s: %w", op, err)
	}

	url := t.googleCalendarService.CreateUrl(ctx, &event)
	message.Answer = "[Google Calendar](" + url + ")"

	file, err := t.appleCalendar.CreateFile(ctx, &event)
	if err != nil {
		t.l.Error(fmt.Errorf("%s: %w", op, err))
		return fmt.Errorf("%s: %w", op, err)
	}
	message.File = file

	return nil
}
