package handler

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"log/slog"

	"gopkg.in/telebot.v3"
)

type button struct {
	lg            *slog.Logger
	tgInvoiceSrv  tgInvoiceSrv
	tgKeyboardSrv tgKeyboardSrv
}

func newButton(
	bot *telebot.Bot,
	lg *slog.Logger,
	tgInvoiceSrv tgInvoiceSrv,
	tgKeyboardSrv tgKeyboardSrv,
) *button {
	b := button{
		lg:            lg,
		tgInvoiceSrv:  tgInvoiceSrv,
		tgKeyboardSrv: tgKeyboardSrv,
	}

	bot.Handle(&telebot.Btn{
		Unique: def.TGBuyMoreQuestionsButton,
	}, b.buyMoreQuestions)
	bot.Handle(&telebot.Btn{
		Unique: def.TGSelectQuestionsCountButton,
	}, b.selectQuestionsCount)

	return &b
}

func (b *button) buyMoreQuestions(c telebot.Context) error {
	const op = "handler.button.buyMoreQuestions"
	b.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := c.Get("ctx").(context.Context)

	if err := c.Delete(); err != nil {
		b.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send("Выберите количество вопросов 🤪", &telebot.ReplyMarkup{
		InlineKeyboard: b.tgKeyboardSrv.Prices(ctx),
	})
}

func (b *button) selectQuestionsCount(c telebot.Context) error {
	const op = "handler.button.selectQuestionsCount"
	b.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := c.Get("ctx").(context.Context)
	user := c.Get("user").(*model.User)

	if err := c.Delete(); err != nil {
		b.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	tgInvoice, err := b.tgInvoiceSrv.CreateBuyInvoice(
		ctx,
		user.ChatID,
		c.Callback().Data,
	)
	if err != nil {
		b.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send(tgInvoice)
}
