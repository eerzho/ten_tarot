package srv

import (
	"bot/internal/model"
	"context"
	"time"
)

type (
	userRepo interface {
		Create(ctx context.Context, user *model.User) error
		Update(ctx context.Context, user *model.User) error
		ExistsByChatID(ctx context.Context, chatID string) bool
		GetByChatID(ctx context.Context, chatID string) (*model.User, error)
	}

	messageRepo interface {
		Create(ctx context.Context, message *model.Message) error
		CountByChatIDAndFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error)
	}

	invoiceRepo interface {
		Create(ctx context.Context, invoice *model.Invoice) error
		Update(ctx context.Context, invoice *model.Invoice) error
		GetByID(ctx context.Context, id string) (*model.Invoice, error)
	}

	supportRequestRepo interface {
		Create(ctx context.Context, supportRequest *model.SupportRequest) error
	}
)
