package handler

import (
	"bot/internal/constant"
	"bot/internal/model"
	"context"
	"time"

	"gopkg.in/telebot.v3"
)

type (
	tgKeyboardService interface {
		Prices(ctx context.Context) [][]telebot.InlineButton
		OverLimit(ctx context.Context) [][]telebot.InlineButton
	}

	tgUserService interface {
		DecreaseQC(ctx context.Context, user *model.TGUser, count int) error
		UpdateState(ctx context.Context, user *model.TGUser, state constant.UserState) error
		GetOrCreateByChatIDUsername(ctx context.Context, chatID, username string) (*model.TGUser, error)
	}

	tgMessageService interface {
		CountByChatIDFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error)
		CreateByChatIDUQ(ctx context.Context, chatID, userQuestion string) (*model.TGMessage, error)
	}

	tgInvoiceService interface {
		IsValidByID(ctx context.Context, id string) bool
		SuccessPayment(ctx context.Context, id, chargeID string, user *model.TGUser) error
		CreateByChatIDData(ctx context.Context, chatID, data string) (*model.TGInvoice, error)
		CreateByChatIDSC(ctx context.Context, chatID string, starsCount int) (*model.TGInvoice, error)
	}

	supportRequestService interface {
		CreateByUserQuestion(ctx context.Context, user *model.TGUser, question string) (*model.SupportRequest, error)
	}
)
