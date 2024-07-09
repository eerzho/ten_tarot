package service

import (
	"context"
	"time"

	"github.com/eerzho/ten_tarot/internal/entity"
)

type (
	TGUserRepo interface {
		Create(ctx context.Context, user *entity.TGUser) error
		ExistsByChatID(ctx context.Context, chatID string) (bool, error)
		Count(ctx context.Context, chatID, username string) (int, error)
		ByChatID(ctx context.Context, chatID string) (*entity.TGUser, error)
		List(ctx context.Context, chatID, username string, page, count int) ([]entity.TGUser, error)
	}

	TGUser struct {
		repo TGUserRepo
	}
)

func NewTGUser(repo TGUserRepo) *TGUser {
	return &TGUser{repo: repo}
}

func (t *TGUser) List(ctx context.Context, username, chatID string, page, count int) ([]entity.TGUser, int, error) {
	users, err := t.repo.List(ctx, username, chatID, page, count)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.count(ctx, chatID, username)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (t *TGUser) Create(ctx context.Context, chatID, username string) (*entity.TGUser, error) {
	exists, _ := t.existsByChatID(ctx, chatID)
	if exists {
		user, err := t.byChatID(ctx, chatID)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	user := entity.TGUser{
		ChatID:    chatID,
		Username:  username,
		CreatedAt: time.Now().Format(time.DateTime),
	}

	if err := t.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (t *TGUser) existsByChatID(ctx context.Context, chatID string) (bool, error) {
	exists, err := t.repo.ExistsByChatID(ctx, chatID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (t *TGUser) byChatID(ctx context.Context, chatID string) (*entity.TGUser, error) {
	user, err := t.repo.ByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (t *TGUser) count(ctx context.Context, chatID, username string) (int, error) {
	count, err := t.repo.Count(ctx, chatID, username)
	if err != nil {
		return 0, err
	}

	return count, nil
}
