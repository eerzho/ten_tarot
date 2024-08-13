package handler

import (
	"bot/internal/model"
	"context"
	"fmt"
	"log/slog"

	"gopkg.in/telebot.v3"
)

type (
	payment struct {
		lg         *slog.Logger
		userSrv    userSrv
		invoiceSrv invoiceSrv
	}
)

func newPayment(
	bot *telebot.Bot,
	lg *slog.Logger,
	userSrv userSrv,
	invoiceSrv invoiceSrv,
) *payment {
	p := payment{
		lg:         lg,
		userSrv:    userSrv,
		invoiceSrv: invoiceSrv,
	}

	bot.Handle(telebot.OnCheckout, p.checkout)
	bot.Handle(telebot.OnPayment, p.payment)

	return &p
}

func (p *payment) checkout(c telebot.Context) error {
	const op = "handler.payment.checkout"
	p.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := c.Get("ctx").(context.Context)

	preCQ := c.PreCheckoutQuery()
	if !p.invoiceSrv.IsValidByID(ctx, preCQ.Payload) {
		return c.Bot().Accept(preCQ, "Вы уже оплатили 🥳")
	} else {
		return c.Bot().Accept(preCQ)
	}
}

func (p *payment) payment(c telebot.Context) error {
	const op = "handler.payment.payment"
	p.lg.Debug(op, slog.Any("RID", c.Get(RID)))
	ctx := c.Get("ctx").(context.Context)
	user := c.Get("user").(*model.User)

	err := p.invoiceSrv.UpdateChargeID(
		ctx,
		c.Message().Payment.Payload,
		c.Message().Payment.TelegramChargeID,
		user,
	)
	if err != nil {
		p.lg.Error(op, slog.String("error", err.Error()))
		return c.Send("✨Пожалуйста, повторите попытку позже✨")
	}

	return c.Send(fmt.Sprintf("У вас %d вопросов 🤯", user.QuestionCount))
}
