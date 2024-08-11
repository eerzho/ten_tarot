package service

import (
	"bot/internal/constant"
	"bot/internal/model"
	"context"
	"log/slog"
)

type (
	TGUser struct {
		lg         *slog.Logger
		tgUserRepo tgUserRepo
	}
)

func NewTGUser(lg *slog.Logger, tgUserRepo tgUserRepo) *TGUser {
	return &TGUser{
		lg:         lg,
		tgUserRepo: tgUserRepo,
	}
}

func (t *TGUser) GetOrCreateByChatIDUsername(ctx context.Context, chatID, username string) (*model.TGUser, error) {
	const op = "service.TGUser.GetOrCreateByChatIDUsername"
	t.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.String("username", username),
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
	t.lg.Debug(
		op,
		slog.Any("user", user),
		slog.Int("count", count),
	)

	user.QuestionCount += count
	if err := t.tgUserRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (t *TGUser) DecreaseQC(ctx context.Context, user *model.TGUser, count int) error {
	const op = "service.TGUser.DecreaseQC"
	t.lg.Debug(
		op,
		slog.Any("user", user),
		slog.Int("count", count),
	)

	user.QuestionCount -= count
	if err := t.tgUserRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (t *TGUser) UpdateState(ctx context.Context, user *model.TGUser, state constant.UserState) error {
	const op = "service.TGUser.UpdateState"
	t.lg.Debug(
		op,
		slog.Any("user", user),
		slog.Any("state", state),
	)

	user.State = state
	if err := t.tgUserRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
