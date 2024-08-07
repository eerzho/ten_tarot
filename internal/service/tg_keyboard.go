package service

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/constant"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type (
	TGKeyboard struct {
	}
)

func NewTGKeyboard() *TGKeyboard {
	return &TGKeyboard{}
}

func (t *TGKeyboard) OverLimit(ctx context.Context) [][]telebot.InlineButton {
	const op = "service.TGKeyboard.OverLimit"
	logger.Debug(op)

	buttons := [][]telebot.InlineButton{
		{
			telebot.InlineButton{
				Unique: constant.BuyMoreQuestionsBTN,
				Text:   "Купите больше вопросов 🤩",
			},
		},
	}

	return buttons
}

func (t *TGKeyboard) Prices(ctx context.Context) [][]telebot.InlineButton {
	const op = "service.TGKeyboard.Prices"
	logger.Debug(op)

	buttons := [][]telebot.InlineButton{
		{
			telebot.InlineButton{
				Unique: constant.SelectQuestionsCountBTN,
				Text:   "5 вопросов - 50 ⭐️",
				Data:   "5:50",
			},
			telebot.InlineButton{
				Unique: constant.SelectQuestionsCountBTN,
				Text:   "10 вопросов - 85 ⭐️",
				Data:   "10:85",
			},
			//telebot.InlineButton{
			//	Unique: constant.SelectQuestionsCount,
			//	Text:   "Test",
			//	Data:   "1:1",
			//},
		},
	}

	return buttons
}
