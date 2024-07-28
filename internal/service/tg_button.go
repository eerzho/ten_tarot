package service

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/constant"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type (
	TGButton struct {
	}
)

func NewTGButton() *TGButton {
	return &TGButton{}
}

func (t *TGButton) OverLimit(ctx context.Context) [][]telebot.InlineButton {
	const op = "service.TGButton.OverLimit"
	logger.Debug(op)

	buttons := [][]telebot.InlineButton{
		{
			telebot.InlineButton{
				Unique: constant.BuyMoreQuestions,
				Text:   "Купите больше вопросов 🤩",
			},
		},
	}

	return buttons
}

func (t *TGButton) Prices(ctx context.Context) [][]telebot.InlineButton {
	const op = "service.TGButton.Prices"
	logger.Debug(op)

	buttons := [][]telebot.InlineButton{
		{
			telebot.InlineButton{
				Unique: constant.SelectQuestionsAmount,
				Text:   "5 вопросов - 50 ⭐️",
				Data:   "5:50",
			},
			telebot.InlineButton{
				Unique: constant.SelectQuestionsAmount,
				Text:   "10 вопросов - 85 ⭐️",
				Data:   "10:85",
			},
			//telebot.InlineButton{
			//	Unique: constant.SelectQuestionsAmount,
			//	Text:   "Test",
			//	Data:   "1:1",
			//},
		},
	}

	return buttons
}
