package service

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/model"
)

type tarotMock struct {
}

func NewTarotMock() Tarot {
	return &tarotMock{}
}

func (t *tarotMock) Oracle(ctx context.Context, question string, hand []model.Card) (string, error) {
	return "From Mock", nil
}
