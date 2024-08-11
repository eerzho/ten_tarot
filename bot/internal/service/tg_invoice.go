package service

import (
	"bot/internal/constant"
	"bot/internal/failure"
	"bot/internal/model"
	"context"
	"log/slog"
	"strconv"
	"strings"
)

type (
	TGInvoice struct {
		lg            *slog.Logger
		tgInvoiceRepo tgInvoiceRepo
		tgUserService tgUserService
	}
)

func NewTGInvoice(
	lg *slog.Logger,
	tgInvoiceRepo tgInvoiceRepo,
	tgUserService tgUserService,
) *TGInvoice {
	return &TGInvoice{
		lg:            lg,
		tgInvoiceRepo: tgInvoiceRepo,
		tgUserService: tgUserService,
	}
}

func (t *TGInvoice) IsValidByID(ctx context.Context, id string) bool {
	const op = "service.TGInvoice.IsValidByID"
	t.lg.Debug(op, slog.String("id", id))

	invoice, err := t.tgInvoiceRepo.GetByID(ctx, id)
	if err != nil {
		return false
	}

	return invoice.ChargeID == ""
}

func (t *TGInvoice) CreateByChatIDData(ctx context.Context, chatID, data string) (*model.TGInvoice, error) {
	const op = "service.TGInvoice.CreateByChatIDData"
	t.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.String("data", data),
	)

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return nil, failure.ErrCallbackData
	}
	count, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, failure.ErrCallbackData
	}
	stars, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, failure.ErrCallbackData
	}

	invoice := model.TGInvoice{
		Type:          constant.InvoiceBuyType,
		ChatID:        chatID,
		StarsCount:    stars,
		QuestionCount: count,
	}

	if err = t.tgInvoiceRepo.Create(ctx, &invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (t *TGInvoice) CreateByChatIDSC(ctx context.Context, chatID string, starsCount int) (*model.TGInvoice, error) {
	const op = "service.TGInvoice.CreateByChatIDSC"
	t.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.Int("starsCount", starsCount),
	)

	invoice := model.TGInvoice{
		Type:          constant.InvoiceDonateType,
		ChatID:        chatID,
		StarsCount:    starsCount,
		QuestionCount: 1,
	}

	if err := t.tgInvoiceRepo.Create(ctx, &invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (t *TGInvoice) SuccessPayment(ctx context.Context, id, chargeID string, user *model.TGUser) error {
	const op = "service.TGInvoice.UpdateByIDChargeID"
	t.lg.Debug(
		op,
		slog.String("id", id),
		slog.String("chargeID", chargeID),
		slog.Any("user", user),
	)

	invoice, err := t.tgInvoiceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	invoice.ChargeID = chargeID
	if err = t.tgInvoiceRepo.Update(ctx, invoice); err != nil {
		return err
	}

	if err = t.tgUserService.IncreaseQC(ctx, user, invoice.QuestionCount); err != nil {
		return err
	}

	return nil
}
