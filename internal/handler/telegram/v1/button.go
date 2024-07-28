package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/constant"
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

	if err := ctx.Delete(); err != nil {
		logger.OPError(op, err)
		return err
	}

	opt := telebot.ReplyMarkup{
		InlineKeyboard: b.tgButtonService.Prices(context.Background()),
	}

	if err := ctx.Send("Выберите количество вопросов 🤪", &opt); err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}

func (b *button) selectQuestionsAmount(ctx telebot.Context) error {
	const op = "handler.telegram.v1.button.selectQuestionsAmount"

	if err := ctx.Delete(); err != nil {
		logger.OPError(op, err)
		if err = ctx.Send("✨Пожалуйста, повторите попытку позже✨"); err != nil {
			logger.OPError(op, err)
			return err
		}
		return err
	}

	invoice, err := b.tgInvoiceService.CreateByChatIDData(
		context.Background(),
		strconv.Itoa(int(ctx.Sender().ID)),
		ctx.Callback().Data,
	)
	if err != nil {
		logger.OPError(op, err)
		if err = ctx.Send("✨Пожалуйста, повторите попытку позже✨"); err != nil {
			logger.OPError(op, err)
			return err
		}
		return err
	}

	in := telebot.Invoice{
		Title: fmt.Sprintf(
			"%d - вопросов",
			invoice.QuestionCount,
		),
		Description: fmt.Sprintf(
			"Вы сможете задать еще %d вопросов",
			invoice.QuestionCount,
		),
		Payload:  invoice.ID,
		Currency: "XTR",
		Prices: []telebot.Price{
			{
				Label: fmt.Sprintf(
					"%d - вопросов",
					invoice.QuestionCount,
				),
				Amount: invoice.StarsCount,
			},
		},
	}

	if err = ctx.Send(&in); err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}
