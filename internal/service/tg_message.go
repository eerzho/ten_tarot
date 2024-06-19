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
		cardService   *Card
		tarotService  *Tarot
	}
)

func NewTGMessage(
	l logger.Logger,
	repo TGMessageRepo,
	tgUserService *TGUser,
	cardService *Card,
	tarotService *Tarot,
) *TGMessage {
	return &TGMessage{
		l:             l,
		repo:          repo,
		tgUserService: tgUserService,
		cardService:   cardService,
		tarotService:  tarotService,
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

	hand := t.cardService.Shuffle(ctx, 5)
	answer, err := t.tarotService.Oracle(ctx, message.Text, hand)
	if err != nil {
		t.l.Error(fmt.Errorf("%s: %w", op, err))
		return fmt.Errorf("%s: %w", op, err)
	}
	message.Answer = answer

	return nil
}
