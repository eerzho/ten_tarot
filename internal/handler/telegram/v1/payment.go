package v1

import (
	"context"
	"fmt"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
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
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}

	preCQ := ctx.PreCheckoutQuery()
	if !p.tgInvoiceService.IsValidByID(oc, preCQ.Payload) {
		return ctx.Bot().Accept(preCQ, "Вы уже оплатили 🥳")
	} else {
		return ctx.Bot().Accept(preCQ)
	}
}

func (p *payment) payment(ctx telebot.Context) error {
	const op = "handler.telegram.v1.payment.payment"
	logger.Debug(op, logger.Any("RID", ctx.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := ctx.Get("user").(*model.TGUser)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}
	oc, ok := ctx.Get("oc").(context.Context)
	if !ok {
		logger.OPError(op, failure.ErrContextData)
		return ctx.Send(errTGMsg)
	}

	err := p.tgInvoiceService.SuccessPayment(
		oc,
		ctx.Message().Payment.Payload,
		ctx.Message().Payment.TelegramChargeID,
		user,
	)
	if err != nil {
		logger.OPError(op, err)
		return ctx.Send(errTGMsg)
	}

	return ctx.Send(fmt.Sprintf("У вас %d вопросов 🤯", user.QuestionCount))
}
