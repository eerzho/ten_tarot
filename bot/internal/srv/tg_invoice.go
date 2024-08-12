package srv

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"gopkg.in/telebot.v3"
)

type TGInvoice struct {
	lg         *slog.Logger
	invoiceSrv invoiceSrv
	userSrv    userSrv
}

func NewTGInvoice(
	lg *slog.Logger,
	invoiceSrv invoiceSrv,
	userSrv userSrv,
) *TGInvoice {
	return &TGInvoice{
		lg:         lg,
		invoiceSrv: invoiceSrv,
		userSrv:    userSrv,
	}
}

func (ti *TGInvoice) CreateBuyInvoice(ctx context.Context, chatID, data string) (*telebot.Invoice, error) {
	const op = "srv.TGInvoice.CreateBuyInvoice"
	ti.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.String("data", data),
	)

	invoice, err := ti.invoiceSrv.CreateByChatIDAndData(ctx, chatID, data)
	if err != nil {
		return nil, err
	}

	tgInvoice := telebot.Invoice{
		Title:       fmt.Sprintf("%d - вопросов", invoice.QuestionCount),
		Description: fmt.Sprintf("Вы сможете задать еще %d вопросов", invoice.QuestionCount),
		Payload:     invoice.ID.String(),
		Currency:    "XTR",
		Prices: []telebot.Price{
			{
				Label:  fmt.Sprintf("%d - вопросов", invoice.QuestionCount),
				Amount: invoice.StarsCount,
			},
		},
	}

	return &tgInvoice, nil
}

func (ti *TGInvoice) CreateDonateInvoice(ctx context.Context, user *model.User, text string) (*telebot.Invoice, error) {
	const op = "srv.TGInvoice.CreateDonateInvoice"
	ti.lg.Debug(
		op,
		slog.Any("user", user),
		slog.String("text", text),
	)

	starsCount, err := strconv.Atoi(text)
	if err != nil {
		return nil, def.ErrInvalidType
	}

	invoice, err := ti.invoiceSrv.CreateByChatIDAndStartsCount(ctx, user.ChatID, starsCount)
	if err != nil {
		return nil, err
	}

	tgInvoice := telebot.Invoice{
		Title:       "Благодарим за поддержку!",
		Description: "Ваш вклад поможет развивать проект и продвигать его дальше!",
		Payload:     invoice.ID.String(),
		Currency:    "XTR",
		Prices: []telebot.Price{
			{
				Label:  "Поддержка проекта",
				Amount: invoice.StarsCount,
			},
		},
	}

	err = ti.userSrv.UpdateState(ctx, user, def.UserDefaultState)
	if err != nil {
		return nil, err
	}

	return &tgInvoice, nil
}
