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
		tgButtonService  tgButtonService
		tgInvoiceService tgInvoiceService
	}
)

func newButton(
	bot *telebot.Bot,
	tgButtonService tgButtonService,
	tgInvoiceService tgInvoiceService,
) *button {
	b := button{
		tgButtonService:  tgButtonService,
		tgInvoiceService: tgInvoiceService,
	}

	bot.Handle(&telebot.Btn{
		Unique: constant.BuyMoreQuestions,
	}, b.buyMoreQuestions)
	bot.Handle(&telebot.Btn{
		Unique: constant.SelectQuestionsAmount,
	}, b.selectQuestionsAmount)

	return &b
}

func (b *button) buyMoreQuestions(ctx telebot.Context) error {
	const op = "handler.telegram.v1.button.buyMoreQuestions"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		if err := ctx.Send(errTGMsg); err != nil {
			logger.OPError(op, err)
			return err
		}
		return failure.ErrContextData
	}

	if err := ctx.Delete(); err != nil {
		logger.OPError(op, err)
		return err
	}

	opt := telebot.ReplyMarkup{
		InlineKeyboard: b.tgButtonService.Prices(oc),
	}
	if err := ctx.Send("Выберите количество вопросов 🤪", &opt); err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}

func (b *button) selectQuestionsAmount(ctx telebot.Context) error {
	const op = "handler.telegram.v1.button.selectQuestionsAmount"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		if err := ctx.Send(errTGMsg); err != nil {
			logger.OPError(op, err)
			return err
		}
		return failure.ErrContextData
	}

	if err := ctx.Delete(); err != nil {
		logger.OPError(op, err)
		if err = ctx.Send(errTGMsg); err != nil {
			logger.OPError(op, err)
		}
		return err
	}

	tgInvoice, err := b.tgInvoiceService.CreateByChatIDData(
		oc,
		strconv.Itoa(int(ctx.Sender().ID)),
		ctx.Callback().Data,
	)
	if err != nil {
		logger.OPError(op, err)
		if err = ctx.Send(errTGMsg); err != nil {
			logger.OPError(op, err)
		}
		return err
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

	if err = ctx.Send(&invoice); err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}
