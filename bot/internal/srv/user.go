package srv

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"log/slog"
)

type User struct {
	lg       *slog.Logger
	userRepo userRepo
}

func NewUser(
	lg *slog.Logger,
	userRepo userRepo,
) *User {
	return &User{
		lg:       lg,
		userRepo: userRepo,
	}
}

func (u *User) GetOrCreateByChatIDAndUsername(ctx context.Context, chatID, username string) (*model.User, error) {
	const op = "srv.User.GetOrCreateByChatIDAndUsername"
	u.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.String("username", username),
	)

	exists := u.userRepo.ExistsByChatID(ctx, chatID)
	if exists {
		user, err := u.userRepo.GetByChatID(ctx, chatID)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	user := model.User{
		ChatID:   chatID,
		Username: username,
	}

	if err := u.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) IncreaseQuestionCount(ctx context.Context, user *model.User, count int) error {
	const op = "srv.User.IncreaseQC"
	u.lg.Debug(
		op,
		slog.Any("user", user),
		slog.Int("count", count),
	)

	user.QuestionCount += count
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *User) DecreaseQuestionCount(ctx context.Context, user *model.User, count int) error {
	const op = "srv.User.DecreaseQC"
	u.lg.Debug(
		op,
		slog.Any("user", user),
		slog.Int("count", count),
	)

	user.QuestionCount -= count
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateState(ctx context.Context, user *model.User, state def.UserState) error {
	const op = "srv.User.UpdateState"
	u.lg.Debug(
		op,
		slog.Any("user", user),
		slog.Any("state", state),
	)

	user.State = state
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
