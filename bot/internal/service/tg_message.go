package service

import (
	"bot/internal/model"
	"context"
	"log/slog"
	"time"
)

type (
	TGMessage struct {
		lg            *slog.Logger
		tgMessageRepo tgMessageRepo
		deckService   deckService
		tarotService  tarotService
	}
)

func NewTGMessage(
	lg *slog.Logger,
	tgMessageRepo tgMessageRepo,
	deckService deckService,
	tarotService tarotService,
) *TGMessage {
	return &TGMessage{
		lg:            lg,
		tgMessageRepo: tgMessageRepo,
		deckService:   deckService,
		tarotService:  tarotService,
	}
}

func (t *TGMessage) CountByChatIDFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error) {
	const op = "service.TGMessage.CountByChatIDFromTime"
	t.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.Any("fromTime", fromTime),
	)

	count, err := t.tgMessageRepo.CountByChatIDFromTime(ctx, chatID, fromTime)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (t *TGMessage) CreateByChatIDUQ(ctx context.Context, chatID, userQuestion string) (*model.TGMessage, error) {
	const op = "service.TGMessage.CreateByChatIDUQ"
	t.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.String("userQuestion", userQuestion),
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
