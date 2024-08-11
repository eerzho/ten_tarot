package service

import (
	"bot/internal/model"
	"context"
	"log/slog"
)

type TarotMock struct {
	lg *slog.Logger
}

func NewTarotMock(lg *slog.Logger) *TarotMock {
	return &TarotMock{
		lg: lg,
	}
}

func (t *TarotMock) Oracle(ctx context.Context, userQuestion string, drawnCards []model.Card) (string, error) {
	const op = "service.TarotMock.Oracle"
	t.lg.Debug(
		op,
		slog.String("userQuestion", userQuestion),
		slog.Any("drawnCardsCount", len(drawnCards)),
	)

	return "From Mock", nil
}
