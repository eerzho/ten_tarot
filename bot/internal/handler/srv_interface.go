package handler

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"time"

	"gopkg.in/telebot.v3"
)

type (
	tgKeyboardSrv interface {
		Prices(ctx context.Context) [][]telebot.InlineButton
		OverLimit(ctx context.Context) [][]telebot.InlineButton
	}

	tgCommandSrv interface {
		GetCommands(ctx context.Context) []telebot.Command
	}

	tgInvoiceSrv interface {
		CreateBuyInvoice(ctx context.Context, chatID, data string) (*telebot.Invoice, error)
		CreateDonateInvoice(ctx context.Context, user *model.User, text string) (*telebot.Invoice, error)
	}

	userSrv interface {
		UpdateState(ctx context.Context, user *model.User, state def.UserState) error
		DecreaseQuestionCount(ctx context.Context, user *model.User, count int) error
		GetOrCreateByChatIDAndUsername(ctx context.Context, chatID, username string) (*model.User, error)
	}

	messageSrv interface {
		CountByChatIDAndFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error)
		CreateByChatIDAndUserQuestion(ctx context.Context, chatID, userQuestion string) (*model.Message, error)
	}

	invoiceSrv interface {
		IsValidByID(ctx context.Context, id string) bool
		UpdateChargeID(ctx context.Context, id, chargeID string, user *model.User) error
	}

	supportRequestSrv interface {
		CreateByUserQuestion(ctx context.Context, user *model.User, question string) (*model.SupportRequest, error)
	}
)
