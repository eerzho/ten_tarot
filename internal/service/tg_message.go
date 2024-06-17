package service

import (
	"context"
	"fmt"

	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/ten_tarot/internal/entity"
)

type (
	TGMessageRepo interface {
		All(ctx context.Context, chatID string, page, count int) ([]entity.TGMessage, error)
		Create(ctx context.Context, message *entity.TGMessage) error
	}

	TGMessage struct {
		l             logger.Logger
		repo          TGMessageRepo
		tgUserService *TGUser
	}
)

func NewTGMessage(
	l logger.Logger,
	repo TGMessageRepo,
	tgUserService *TGUser,
) *TGMessage {
	return &TGMessage{
		l:             l,
		repo:          repo,
		tgUserService: tgUserService,
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

	return nil
}
