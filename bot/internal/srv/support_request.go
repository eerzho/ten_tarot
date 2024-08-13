package srv

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"log/slog"
)

type SupportRequest struct {
	lg                 *slog.Logger
	supportRequestRepo supportRequestRepo
	userSrv            userSrv
}

func NewSupportRequest(
	lg *slog.Logger,
	supportRequestRepo supportRequestRepo,
	userSrv userSrv,
) *SupportRequest {
	return &SupportRequest{
		lg:                 lg,
		supportRequestRepo: supportRequestRepo,
		userSrv:            userSrv,
	}
}

func (s *SupportRequest) CreateByUserQuestion(ctx context.Context, user *model.User, question string) (*model.SupportRequest, error) {
	const op = "srv.SupportRequest.CreateByUserQuestion"
	s.lg.Debug(
		op,
		slog.Any("user", user),
		slog.String("question", question),
	)

	sp := model.SupportRequest{
		ChatID:   user.ChatID,
		Question: question,
	}

	if err := s.supportRequestRepo.Create(ctx, &sp); err != nil {
		return nil, err
	}

	if err := s.userSrv.UpdateState(ctx, user, def.UserDefaultState); err != nil {
		return nil, err
	}

	return &sp, nil
}
