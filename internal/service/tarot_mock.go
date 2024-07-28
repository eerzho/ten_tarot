package service

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
)

type TarotMock struct {
}

func NewTarotMock() *TarotMock {
	return &TarotMock{}
}

func (t *TarotMock) Oracle(ctx context.Context, userQuestion string, drawnCards []model.Card) (string, error) {
	const op = "service.TarotMock.Oracle"
	logger.Debug(
		op,
		logger.Any("userQuestion", userQuestion),
		logger.Any("drawnCardsCount", len(drawnCards)),
	)
	return "From Mock", nil
}
