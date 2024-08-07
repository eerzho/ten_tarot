package service

import (
	"context"
	"time"

	"github.com/eerzho/ten_tarot/internal/model"
)

type (
	tgUserRepo interface {
		Create(ctx context.Context, user *model.TGUser) error
		Update(ctx context.Context, user *model.TGUser) error
		ExistsByChatID(ctx context.Context, chatID string) bool
		GetByChatID(ctx context.Context, chatID string) (*model.TGUser, error)
		GetListCount(ctx context.Context, chatID, username string) (int, error)
		GetList(ctx context.Context, username, chatID string, page, count int) ([]model.TGUser, error)
	}

	tgMessageRepo interface {
		Create(ctx context.Context, message *model.TGMessage) error
		GetListCount(ctx context.Context, chatID string) (int, error)
		CountByChatIDFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error)
		GetList(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, error)
	}

	tgInvoiceRepo interface {
		Create(ctx context.Context, invoice *model.TGInvoice) error
		Update(ctx context.Context, invoice *model.TGInvoice) error
		GetByID(ctx context.Context, id string) (*model.TGInvoice, error)
	}

	deckService interface {
		Shuffle(ctx context.Context, n int) ([]model.Card, error)
	}

	tarotService interface {
		Oracle(ctx context.Context, question string, hand []model.Card) (string, error)
	}

	tgUserService interface {
		IncreaseQC(ctx context.Context, user *model.TGUser, count int) error
		UpdateState(ctx context.Context, user *model.TGUser, state string) error
	}

	supportRequestRepo interface {
		Create(ctx context.Context, supportRequest *model.SupportRequest) error
	}
)
