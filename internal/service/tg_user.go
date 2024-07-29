package service

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
)

type (
	TGUser struct {
		tgUserRepo tgUserRepo
	}
)

func NewTGUser(tgUserRepo tgUserRepo) *TGUser {
	return &TGUser{
		tgUserRepo: tgUserRepo,
	}
}

func (t *TGUser) GetList(ctx context.Context, username, chatID string, page, count int) ([]model.TGUser, int, error) {
	const op = "service.TGUser.GetList"
	logger.Debug(
		op,
		logger.Any("username", username),
		logger.Any("chatID", chatID),
		logger.Any("page", page),
		logger.Any("count", count),
	)

	users, err := t.tgUserRepo.GetList(
		ctx,
		username,
		chatID,
		page,
		count,
	)
	if err != nil {
		return nil, 0, err
	}

	total, err := t.tgUserRepo.GetListCount(
		ctx,
		chatID,
		username,
	)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (t *TGUser) GetOrCreateByChatIDUsername(ctx context.Context, chatID, username string) (*model.TGUser, error) {
	const op = "service.TGUser.GetOrCreateByChatIDUsername"
	logger.Debug(op,
		logger.Any("chatID", chatID),
		logger.Any("username", username),
	)

	exists := t.tgUserRepo.ExistsByChatID(ctx, chatID)
	if exists {
		user, err := t.tgUserRepo.GetByChatID(ctx, chatID)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	user := model.TGUser{
		ChatID:   chatID,
		Username: username,
	}

	if err := t.tgUserRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (t *TGUser) IncreaseQC(ctx context.Context, user *model.TGUser, count int) error {
	const op = "service.TGUser.IncreaseQC"
	logger.Debug(
		op,
		logger.Any("user", user),
		logger.Any("count", count),
	)

	user.QuestionCount += count
	if err := t.tgUserRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (t *TGUser) DecreaseQC(ctx context.Context, user *model.TGUser, count int) error {
	const op = "service.TGUser.DecreaseQC"
	logger.Debug(
		op,
		logger.Any("user", user),
		logger.Any("count", count),
	)

	user.QuestionCount -= count
	if err := t.tgUserRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
