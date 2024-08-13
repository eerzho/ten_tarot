package srv

import (
	"bot/internal/dto"
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

func (t *TarotMock) Oracle(ctx context.Context, userQuestion string, drawnCards []dto.Card) (string, error) {
	const op = "srv.TarotMock.Oracle"
	t.lg.Debug(
		op,
		slog.String("userQuestion", userQuestion),
		slog.Any("drawnCardsCount", len(drawnCards)),
	)

	return "From Mock", nil
}
