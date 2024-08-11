package service

import (
	"bot/internal/constant"
	"context"
	"gopkg.in/telebot.v3"
	"log/slog"
)

type (
	TGKeyboard struct {
		lg *slog.Logger
	}
)

func NewTGKeyboard(lg *slog.Logger) *TGKeyboard {
	return &TGKeyboard{
		lg: lg,
	}
}

func (t *TGKeyboard) OverLimit(ctx context.Context) [][]telebot.InlineButton {
	const op = "service.TGKeyboard.OverLimit"
	t.lg.Debug(op)

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
	t.lg.Debug(op)

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
		},
	}

	return buttons
}
