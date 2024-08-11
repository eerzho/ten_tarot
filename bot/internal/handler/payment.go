package handler

import (
	"bot/internal/failure"
	"bot/internal/model"
	"context"
	"fmt"
	"log/slog"

	"gopkg.in/telebot.v3"
)

type (
	payment struct {
		lg               *slog.Logger
		tgInvoiceService tgInvoiceService
		tgUserService    tgUserService
	}
)

func newPayment(
	bot *telebot.Bot,
	lg *slog.Logger,
	tgInvoiceService tgInvoiceService,
	tgUserService tgUserService,
) *payment {
	p := payment{
		lg:               lg,
		tgInvoiceService: tgInvoiceService,
		tgUserService:    tgUserService,
	}

	bot.Handle(telebot.OnCheckout, p.checkout)
	bot.Handle(telebot.OnPayment, p.payment)

	return &p
}

func (p *payment) checkout(c telebot.Context) error {
	const op = "handler.payment.checkout"
	p.lg.Debug(op, slog.Any("RID", c.Get(RID)))

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	oc, ok := c.Get("oc").(context.Context)
	if !ok {
		p.lg.Error(op, slog.String("error", failure.ErrContextData.Error()))
		return c.Send(errTGMsg)
	}

	preCQ := c.PreCheckoutQuery()
	if !p.tgInvoiceService.IsValidByID(oc, preCQ.Payload) {
		return c.Bot().Accept(preCQ, "Вы уже оплатили 🥳")
	} else {
		return c.Bot().Accept(preCQ)
	}
}

func (p *payment) payment(c telebot.Context) error {
	const op = "handler.payment.payment"
	p.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := context.Background()

	errTGMsg := "✨Пожалуйста, повторите попытку позже✨"

	user, ok := c.Get("user").(*model.TGUser)
	if !ok {
		p.lg.Error(op, slog.String("error", failure.ErrContextData.Error()))
		return c.Send(errTGMsg)
	}

	err := p.tgInvoiceService.SuccessPayment(
		ctx,
		c.Message().Payment.Payload,
		c.Message().Payment.TelegramChargeID,
		user,
	)
	if err != nil {
		p.lg.Error(op, slog.String("error", err.Error()))
		return c.Send(errTGMsg)
	}

	return c.Send(fmt.Sprintf("У вас %d вопросов 🤯", user.QuestionCount))
}
