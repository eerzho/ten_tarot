package v1

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/model"
)

type (
	tgUserService interface {
		GetList(ctx context.Context, username, chatID string, page, count int) ([]model.TGUser, int, error)
	}

	tgMessageService interface {
		GetList(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, int, error)
	}
)
