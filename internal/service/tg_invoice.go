package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
)

type (
	TGInvoice struct {
		tgInvoiceRepo tgInvoiceRepo
	}

	tgInvoiceRepo interface {
		Create(ctx context.Context, invoice *model.TGInvoice) error
		Update(ctx context.Context, invoice *model.TGInvoice) error
		GetByID(ctx context.Context, id string) (*model.TGInvoice, error)
	}
)

func NewTGInvoice(tgInvoiceRepo tgInvoiceRepo) *TGInvoice {
	return &TGInvoice{
		tgInvoiceRepo: tgInvoiceRepo,
	}
}

func (t *TGInvoice) CreateByData(ctx context.Context, chatID, data string) (*model.TGInvoice, error) {
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
		Stars:         stars,
		ChatID:        chatID,
		QuestionCount: count,
	}

	if err = t.tgInvoiceRepo.Create(ctx, &invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (t *TGInvoice) IsValidByID(ctx context.Context, id string) bool {
	invoice, err := t.getByID(ctx, id)
	if err != nil {
		return false
	}

	return invoice.ChargeID == ""
}

func (t *TGInvoice) UpdateByIDChargeID(ctx context.Context, id, chargeID string) (*model.TGInvoice, error) {
	invoice, err := t.getByID(ctx, id)
	if err != nil {
		return nil, err
	}

	invoice.ChargeID = chargeID

	if err = t.tgInvoiceRepo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (t *TGInvoice) getByID(ctx context.Context, id string) (*model.TGInvoice, error) {
	invoice, err := t.tgInvoiceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}
