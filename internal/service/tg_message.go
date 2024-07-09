package service

import (
	"context"
	"time"

	"github.com/eerzho/ten_tarot/internal/model"
)

type (
	TGMessageRepo interface {
		Count(ctx context.Context, chatID string) (int, error)
		Create(ctx context.Context, message *model.TGMessage) error
		CountByTime(ctx context.Context, chatID string, st time.Time) (int, error)
		List(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, error)
	}

	TGMessage interface {
		CountByTime(ctx context.Context, chatID string, st time.Time) (int, error)
		Create(ctx context.Context, chatID, text string) (*model.TGMessage, error)
		List(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, int, error)
	}

	tgMessage struct {
		repo         TGMessageRepo
		cardService  Card
		tarotService Tarot
	}
)

func NewTGMessage(
	repo TGMessageRepo,
	cardService Card,
	tarotService Tarot,
) TGMessage {
	return &tgMessage{
		repo:         repo,
		cardService:  cardService,
		tarotService: tarotService,
	}
}

func (t *tgMessage) CountByTime(ctx context.Context, chatID string, st time.Time) (int, error) {
	count, err := t.repo.CountByTime(ctx, chatID, st)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (t *tgMessage) Create(ctx context.Context, chatID, text string) (*model.TGMessage, error) {
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

func (t *tgMessage) List(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, int, error) {
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

func (t *tgMessage) count(ctx context.Context, chatID string) (int, error) {
	count, err := t.repo.Count(ctx, chatID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
