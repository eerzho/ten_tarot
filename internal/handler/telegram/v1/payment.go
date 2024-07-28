package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type (
	payment struct {
		tgInvoiceService tgInvoiceService
		tgUserService    tgUserService
	}
)

func newPayment(
	bot *telebot.Bot,
	tgInvoiceService tgInvoiceService,
	tgUserService tgUserService,
) *payment {
	p := payment{
		tgInvoiceService: tgInvoiceService,
		tgUserService:    tgUserService,
	}

	bot.Handle(telebot.OnCheckout, p.checkout)
	bot.Handle(telebot.OnPayment, p.payment)

	return &p
}

func (p *payment) checkout(ctx telebot.Context) error {
	const op = "handler.telegram.v1.payment.checkout"

	var err error
	preCQ := ctx.PreCheckoutQuery()
	if !p.tgInvoiceService.IsValidByID(
		context.Background(),
		preCQ.Payload,
	) {
		err = ctx.Bot().Accept(preCQ, "Вы уже оплатили 🥳")
	} else {
		err = ctx.Bot().Accept(preCQ)
	}

	if err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}

func (p *payment) payment(ctx telebot.Context) error {
	const op = "handler.telegram.v1.payment.payment"

	invoice, err := p.tgInvoiceService.UpdateByIDChargeID(
		context.Background(),
		ctx.Message().Payment.Payload,
		ctx.Message().Payment.TelegramChargeID,
	)
	if err != nil {
		logger.OPError(op, err)
		return err
	}

	user, err := p.tgUserService.UpdateByChatIDQC(
		context.Background(),
		strconv.Itoa(int(ctx.Sender().ID)), invoice.QuestionCount,
	)
	if err != nil {
		logger.OPError(op, err)
		return err
	}

	if err = ctx.Send(fmt.Sprintf("У вас %d вопросов 🤯", user.QuestionCount)); err != nil {
		logger.OPError(op, err)
		return err
	}

	return nil
}
