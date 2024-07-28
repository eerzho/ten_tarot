package service

import (
	"context"
	"time"

	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
)

type (
	TGMessage struct {
		tgMessageRepo tgMessageRepo
		deckService   deckService
		tarotService  tarotService
	}
)

func NewTGMessage(
	tgMessageRepo tgMessageRepo,
	deckService deckService,
	tarotService tarotService,
) *TGMessage {
	return &TGMessage{
		tgMessageRepo: tgMessageRepo,
		deckService:   deckService,
		tarotService:  tarotService,
	}
}

func (t *TGMessage) CountByChatIDFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error) {
	const op = "service.TGMessage.CountByChatIDFromTime"
	logger.Debug(
		op,
		logger.Any("chatID", chatID),
		logger.Any("fromTime", fromTime),
	)

	count, err := t.tgMessageRepo.CountByChatIDFromTime(ctx, chatID, fromTime)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (t *TGMessage) GetList(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, int, error) {
	const op = "service.TGMessage.GetList"
	logger.Debug(
		op,
		logger.Any("chatID", chatID),
		logger.Any("page", page),
		logger.Any("count", count),
	)

	messages, err := t.tgMessageRepo.GetList(
		ctx,
		chatID,
		page,
		count,
	)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.tgMessageRepo.GetListCount(ctx, chatID)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (t *TGMessage) CreateByChatIDUQ(ctx context.Context, chatID, userQuestion string) (*model.TGMessage, error) {
	const op = "service.TGMessage.CreateByChatIDUQ"
	logger.Debug(
		op,
		logger.Any("chatID", chatID),
		logger.Any("userQuestion", userQuestion),
	)

	hand, err := t.deckService.Shuffle(ctx, 5)
	if err != nil {
		return nil, err
	}

	botAnswer, err := t.tarotService.Oracle(ctx, userQuestion, hand)
	if err != nil {
		return nil, err
	}

	message := model.TGMessage{
		ChatID:       chatID,
		BotAnswer:    botAnswer,
		UserQuestion: userQuestion,
	}
	if err = t.tgMessageRepo.Create(ctx, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
