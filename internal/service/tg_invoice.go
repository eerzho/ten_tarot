package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
)

type (
	TGInvoice struct {
		tgInvoiceRepo tgInvoiceRepo
	}
)

func NewTGInvoice(tgInvoiceRepo tgInvoiceRepo) *TGInvoice {
	return &TGInvoice{
		tgInvoiceRepo: tgInvoiceRepo,
	}
}

func (t *TGInvoice) IsValidByID(ctx context.Context, id string) bool {
	const op = "service.TGInvoice.IsValidByID"
	logger.Debug(op, logger.Any("id", id))

	invoice, err := t.tgInvoiceRepo.GetByID(ctx, id)
	if err != nil {
		return false
	}

	return invoice.ChargeID == ""
}

func (t *TGInvoice) GetByID(ctx context.Context, id string) (*model.TGInvoice, error) {
	const op = "service.TGInvoice.GetByID"
	logger.Debug(op, logger.Any("id", id))

	invoice, err := t.tgInvoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (t *TGInvoice) CreateByChatIDData(ctx context.Context, chatID, data string) (*model.TGInvoice, error) {
	const op = "service.TGInvoice.CreateByChatIDData"
	logger.Debug(
		op,
		logger.Any("chatID", chatID),
		logger.Any("data", data),
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
		ChatID:        chatID,
		StarsCount:    stars,
		QuestionCount: count,
	}

	if err = t.tgInvoiceRepo.Create(ctx, &invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (t *TGInvoice) UpdateByIDChargeID(ctx context.Context, id, chargeID string) (*model.TGInvoice, error) {
	const op = "service.TGInvoice.UpdateByIDChargeID"
	logger.Debug(
		op,
		logger.Any("id", id),
		logger.Any("chargeID", chargeID),
	)

	invoice, err := t.tgInvoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	invoice.ChargeID = chargeID

	if err = t.tgInvoiceRepo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}
