package srv

import (
	"bot/internal/def"
	"bot/internal/dto"
	"bot/internal/model"
	"context"
)

type (
	deckSrv interface {
		Shuffle(ctx context.Context, n int) ([]dto.Card, error)
	}

	tarotSrv interface {
		Oracle(ctx context.Context, question string, hand []dto.Card) (string, error)
	}

	userSrv interface {
		IncreaseQuestionCount(ctx context.Context, user *model.User, count int) error
		UpdateState(ctx context.Context, user *model.User, state def.UserState) error
	}

	invoiceSrv interface {
		CreateByChatIDAndData(ctx context.Context, chatID, data string) (*model.Invoice, error)
		CreateByChatIDAndStartsCount(ctx context.Context, chatID string, starsCount int) (*model.Invoice, error)
	}
)
