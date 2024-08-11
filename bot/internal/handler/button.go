package handler

import (
	"bot/internal/constant"
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
	"log/slog"
	"strconv"
)

type (
	button struct {
		lg                *slog.Logger
		tgKeyboardService tgKeyboardService
		tgInvoiceService  tgInvoiceService
	}
)

func newButton(
	bot *telebot.Bot,
	lg *slog.Logger,
	tgKeyboardService tgKeyboardService,
	tgInvoiceService tgInvoiceService,
) *button {
	b := button{
		lg:                lg,
		tgKeyboardService: tgKeyboardService,
		tgInvoiceService:  tgInvoiceService,
	}

	bot.Handle(&telebot.Btn{
		Unique: constant.BuyMoreQuestionsBTN,
	}, b.buyMoreQuestions)
	bot.Handle(&telebot.Btn{
		Unique: constant.SelectQuestionsCountBTN,
	}, b.selectQuestionsCount)

	return &b
}

func (b *button) buyMoreQuestions(c telebot.Context) error {
	const op = "handler.button.buyMoreQuestions"
	b.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := context.Background()

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	if err := c.Delete(); err != nil {
		b.lg.Error(op, slog.String("error", err.Error()))
		return c.Send(errTGMsg)
	}

	opt := telebot.ReplyMarkup{
		InlineKeyboard: b.tgKeyboardService.Prices(ctx),
	}

	return c.Send("Выберите количество вопросов 🤪", &opt)
}

func (b *button) selectQuestionsCount(c telebot.Context) error {
	const op = "handler.button.selectQuestionsCount"
	b.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := context.Background()

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	if err := c.Delete(); err != nil {
		b.lg.Error(op, slog.String("error", err.Error()))
		return c.Send(errTGMsg)
	}

	tgInvoice, err := b.tgInvoiceService.CreateByChatIDData(
		ctx,
		strconv.Itoa(int(c.Sender().ID)),
		c.Callback().Data,
	)
	if err != nil {
		b.lg.Error(op, slog.String("error", err.Error()))
		return c.Send(errTGMsg)
	}

	invoice := telebot.Invoice{
		Title: fmt.Sprintf(
			"%d - вопросов",
			tgInvoice.QuestionCount,
		),
		Description: fmt.Sprintf(
			"Вы сможете задать еще %d вопросов",
			tgInvoice.QuestionCount,
		),
		Payload:  tgInvoice.ID,
		Currency: "XTR",
		Prices: []telebot.Price{
			{
				Label: fmt.Sprintf(
					"%d - вопросов",
					tgInvoice.QuestionCount,
				),
				Amount: tgInvoice.StarsCount,
			},
		},
	}

	return c.Send(&invoice)
}
