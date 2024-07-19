package service

import (
	"context"
	"time"

	"github.com/eerzho/ten_tarot/internal/model"
)

type (
	TGMessage struct {
		repo         tgMessageRepo
		cardService  cardService
		tarotService tarotService
	}

	tgMessageRepo interface {
		Count(ctx context.Context, chatID string) (int, error)
		Create(ctx context.Context, message *model.TGMessage) error
		CountByTime(ctx context.Context, chatID string, st time.Time) (int, error)
		List(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, error)
	}

	cardService interface {
		Shuffle(ctx context.Context, n int) ([]model.Card, error)
	}

	tarotService interface {
		Oracle(ctx context.Context, question string, hand []model.Card) (string, error)
	}
)

func NewTGMessage(
	repo tgMessageRepo,
	cardService cardService,
	tarotService tarotService,
) *TGMessage {
	return &TGMessage{
		repo:         repo,
		cardService:  cardService,
		tarotService: tarotService,
	}
}

func (t *TGMessage) CountByTime(ctx context.Context, chatID string, st time.Time) (int, error) {
	count, err := t.repo.CountByTime(ctx, chatID, st)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (t *TGMessage) Create(ctx context.Context, chatID, text string) (*model.TGMessage, error) {
	hand, err := t.cardService.Shuffle(ctx, 5)
	if err != nil {
		return nil, err
	}

	answer, err := t.tarotService.Oracle(ctx, text, hand)
	if err != nil {
		return nil, err
	}

	message := model.TGMessage{
		ChatID:    chatID,
		Text:      text,
		Answer:    answer,
		CreatedAt: time.Now().Format(time.DateTime),
	}
	if err = t.repo.Create(ctx, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (t *TGMessage) List(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, int, error) {
	messages, err := t.repo.List(ctx, chatID, page, count)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.count(ctx, chatID)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (t *TGMessage) count(ctx context.Context, chatID string) (int, error) {
	count, err := t.repo.Count(ctx, chatID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
