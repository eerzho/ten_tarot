package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/ten_tarot/internal/entity"
	"github.com/eerzho/ten_tarot/internal/failure"
)

type (
	TGUserRepo interface {
		All(ctx context.Context, username, chatID string, page, count int) ([]entity.TGUser, error)
		ByChatID(ctx context.Context, chatID string) (*entity.TGUser, error)
		Create(ctx context.Context, user *entity.TGUser) error
	}

	TGUser struct {
		l    logger.Logger
		repo TGUserRepo
	}
)

func NewTGUser(l logger.Logger, repo TGUserRepo) *TGUser {
	return &TGUser{l: l, repo: repo}
}

func (t *TGUser) All(ctx context.Context, username, chatID string, page, count int) ([]entity.TGUser, error) {
	const op = "./internal/service/tg_user::All"

	users, err := t.repo.All(ctx, username, chatID, page, count)
	if err != nil {
		t.l.Debug(fmt.Errorf("%s: %w", op, err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (t *TGUser) ByChatID(ctx context.Context, chatID string) (*entity.TGUser, error) {
	const op = "./internal/service/tg_user::ByChatID"

	user, err := t.repo.ByChatID(ctx, chatID)
	if err != nil {
		t.l.Debug(fmt.Errorf("%s: %w", op, err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (t *TGUser) Create(ctx context.Context, user *entity.TGUser) error {
	const op = "./internal/service/tg_user::Create"

	exUser, err := t.ByChatID(ctx, user.ChatID)
	if err != nil {
		t.l.Debug(fmt.Errorf("%s: %w", op, err))
		if !errors.Is(err, failure.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if exUser != nil {
		user = exUser
		return nil
	}

	if err = t.repo.Create(ctx, user); err != nil {
		t.l.Debug(fmt.Errorf("%s: %w", op, err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
