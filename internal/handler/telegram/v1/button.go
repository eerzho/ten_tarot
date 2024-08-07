package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/constant"
	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type (
	button struct {
		tgKeyboardService tgKeyboardService
		tgInvoiceService  tgInvoiceService
	}
)

func newButton(
	bot *telebot.Bot,
	tgKeyboardService tgKeyboardService,
	tgInvoiceService tgInvoiceService,
) *button {
	b := button{
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

func (b *button) buyMoreQuestions(ctx telebot.Context) error {
	const op = "handler.telegram.v1.button.buyMoreQuestions"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}

	if err := ctx.Delete(); err != nil {
		logger.OPError(op, err)
		return ctx.Send(errTGMsg)
	}

	opt := telebot.ReplyMarkup{
		InlineKeyboard: b.tgKeyboardService.Prices(oc),
	}

	return ctx.Send("Выберите количество вопросов 🤪", &opt)
}

func (b *button) selectQuestionsCount(ctx telebot.Context) error {
	const op = "handler.telegram.v1.button.selectQuestionsCount"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}

	if err := ctx.Delete(); err != nil {
		logger.OPError(op, err)
		return ctx.Send(errTGMsg)
	}

	tgInvoice, err := b.tgInvoiceService.CreateByChatIDData(
		oc,
		strconv.Itoa(int(ctx.Sender().ID)),
		ctx.Callback().Data,
	)
	if err != nil {
		logger.OPError(op, err)
		return ctx.Send(errTGMsg)
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

	return ctx.Send(&invoice)
}
