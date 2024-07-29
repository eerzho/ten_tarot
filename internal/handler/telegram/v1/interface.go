package v1

import (
	"context"
	"time"

	"github.com/eerzho/ten_tarot/internal/model"
	"gopkg.in/telebot.v3"
)

type (
	tgButtonService interface {
		Prices(ctx context.Context) [][]telebot.InlineButton
		OverLimit(ctx context.Context) [][]telebot.InlineButton
	}

	tgUserService interface {
		IncreaseQC(ctx context.Context, user *model.TGUser, count int) error
		DecreaseQC(ctx context.Context, user *model.TGUser, count int) error
		GetByChatID(ctx context.Context, chatID string) (*model.TGUser, error)
		GetOrCreateByChatIDUsername(ctx context.Context, chatID, username string) (*model.TGUser, error)
	}

	tgMessageService interface {
		CountByChatIDFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error)
		CreateByChatIDUQ(ctx context.Context, chatID, userQuestion string) (*model.TGMessage, error)
	}

	tgInvoiceService interface {
		IsValidByID(ctx context.Context, id string) bool
		CreateByChatIDData(ctx context.Context, chatID, data string) (*model.TGInvoice, error)
		UpdateByIDChargeID(ctx context.Context, id, chargeID string) (*model.TGInvoice, error)
	}
)
