package srv

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"log/slog"
	"strconv"
	"strings"
)

type Invoice struct {
	lg          *slog.Logger
	invoiceRepo invoiceRepo
	userSrv     userSrv
}

func NewInvoice(
	lg *slog.Logger,
	invoiceRepo invoiceRepo,
	userSrv userSrv,
) *Invoice {
	return &Invoice{
		lg:          lg,
		invoiceRepo: invoiceRepo,
		userSrv:     userSrv,
	}
}

func (i *Invoice) IsValidByID(ctx context.Context, id string) bool {
	const op = "srv.Invoice.IsValidByID"
	i.lg.Debug(op, slog.String("id", id))

	invoice, err := i.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return false
	}

	return invoice.ChargeID == ""
}

func (i *Invoice) CreateByChatIDAndData(ctx context.Context, chatID, data string) (*model.Invoice, error) {
	const op = "srv.Invoice.CreateByChatIDAndData"
	i.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.String("data", data),
	)

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return nil, def.ErrCallbackData
	}
	count, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, def.ErrCallbackData
	}
	stars, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, def.ErrCallbackData
	}

	invoice := model.Invoice{
		Type:          def.InvoiceBuyType,
		ChatID:        chatID,
		StarsCount:    stars,
		QuestionCount: count,
	}

	if err = i.invoiceRepo.Create(ctx, &invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (i *Invoice) CreateByChatIDAndStartsCount(ctx context.Context, chatID string, starsCount int) (*model.Invoice, error) {
	const op = "srv.Invoice.CreateByChatIDAndStartsCount"
	i.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.Int("starsCount", starsCount),
	)

	invoice := model.Invoice{
		Type:          def.InvoiceDonateType,
		ChatID:        chatID,
		StarsCount:    starsCount,
		QuestionCount: 1,
	}

	if err := i.invoiceRepo.Create(ctx, &invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (i *Invoice) UpdateChargeID(ctx context.Context, id, chargeID string, user *model.User) error {
	const op = "srv.Invoice.UpdateChargeID"
	i.lg.Debug(
		op,
		slog.String("id", id),
		slog.String("chargeID", chargeID),
		slog.Any("user", user),
	)

	invoice, err := i.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	invoice.ChargeID = chargeID
	if err = i.invoiceRepo.Update(ctx, invoice); err != nil {
		return err
	}

	if err = i.userSrv.IncreaseQuestionCount(ctx, user, invoice.QuestionCount); err != nil {
		return err
	}

	return nil
}
