package srv

import (
	"bot/internal/model"
	"context"
	"log/slog"
	"time"
)

type Message struct {
	lg          *slog.Logger
	messageRepo messageRepo
	deckSrv     deckSrv
	tarotSrv    tarotSrv
}

func NewMessage(
	lg *slog.Logger,
	messageRepo messageRepo,
	deckSrv deckSrv,
	tarotSrv tarotSrv,
) *Message {
	return &Message{
		lg:          lg,
		messageRepo: messageRepo,
		deckSrv:     deckSrv,
		tarotSrv:    tarotSrv,
	}
}

func (m *Message) CountByChatIDAndFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error) {
	const op = "srv.Message.CountByChatIDAndFromTime"
	m.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.Any("fromTime", fromTime),
	)

	count, err := m.messageRepo.CountByChatIDAndFromTime(ctx, chatID, fromTime)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *Message) CreateByChatIDAndUserQuestion(ctx context.Context, chatID, userQuestion string) (*model.Message, error) {
	const op = "srv.Message.CreateByChatIDAndUserQuestion"
	m.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.String("userQuestion", userQuestion),
	)

	hand, err := m.deckSrv.Shuffle(ctx, 5)
	if err != nil {
		return nil, err
	}

	botAnswer, err := m.tarotSrv.Oracle(ctx, userQuestion, hand)
	if err != nil {
		return nil, err
	}

	message := model.Message{
		ChatID:       chatID,
		BotAnswer:    botAnswer,
		UserQuestion: userQuestion,
	}
	if err = m.messageRepo.Create(ctx, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
