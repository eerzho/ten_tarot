package service

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/model"
)

type TarotMock struct {
}

func NewTarotMock() *TarotMock {
	return &TarotMock{}
}

func (t *TarotMock) Oracle(ctx context.Context, question string, hand []model.Card) (string, error) {
	return "From Mock", nil
}
