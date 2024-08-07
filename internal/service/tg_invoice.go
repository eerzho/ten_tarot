package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/eerzho/ten_tarot/internal/constant"
	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
)

type (
	TGInvoice struct {
		tgInvoiceRepo tgInvoiceRepo
		tgUserService tgUserService
	}
)

func NewTGInvoice(
	tgInvoiceRepo tgInvoiceRepo,
	tgUserService tgUserService,
) *TGInvoice {
	return &TGInvoice{
		tgInvoiceRepo: tgInvoiceRepo,
		tgUserService: tgUserService,
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
	logger.Debug(
		op,
		logger.Any("chatID", chatID),
		logger.Any("starsCount", starsCount),
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
	logger.Debug(
		op,
		logger.Any("id", id),
		logger.Any("user", user),
		logger.Any("chargeID", chargeID),
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
