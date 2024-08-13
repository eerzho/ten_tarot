package srv

import (
	"bot/internal/def"
	"context"
	"log/slog"

	"gopkg.in/telebot.v3"
)

type TGKeyboard struct {
	lg *slog.Logger
}

func NewTGKeyboard(lg *slog.Logger) *TGKeyboard {
	return &TGKeyboard{
		lg: lg,
	}
}

func (t *TGKeyboard) OverLimit(ctx context.Context) [][]telebot.InlineButton {
	const op = "srv.TGKeyboard.OverLimit"
	t.lg.Debug(op)

	buttons := [][]telebot.InlineButton{
		{
			telebot.InlineButton{
				Unique: def.TGBuyMoreQuestionsButton,
				Text:   "Купите больше вопросов 🤩",
			},
		},
	}

	return buttons
}

func (t *TGKeyboard) Prices(ctx context.Context) [][]telebot.InlineButton {
	const op = "srv.TGKeyboard.Prices"
	t.lg.Debug(op)

	buttons := [][]telebot.InlineButton{
		{
			telebot.InlineButton{
				Unique: def.TGSelectQuestionsCountButton,
				Text:   "5 вопросов - 50 ⭐️",
				Data:   "5:50",
			},
			telebot.InlineButton{
				Unique: def.TGSelectQuestionsCountButton,
				Text:   "10 вопросов - 85 ⭐️",
				Data:   "10:85",
			},
		},
	}

	return buttons
}
