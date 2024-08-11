package service

import (
	"bot/internal/constant"
	"bot/internal/model"
	"context"
	"log/slog"
)

type SupportRequest struct {
	lg                 *slog.Logger
	supportRequestRepo supportRequestRepo
	tgUserService      tgUserService
}

func NewSupportRequest(
	lg *slog.Logger,
	supportRequestRepo supportRequestRepo,
	tgUserService tgUserService,
) *SupportRequest {
	return &SupportRequest{
		lg:                 lg,
		supportRequestRepo: supportRequestRepo,
		tgUserService:      tgUserService,
	}
}

func (s *SupportRequest) CreateByUserQuestion(ctx context.Context, user *model.TGUser, question string) (*model.SupportRequest, error) {
	const op = "service.SupportRequest.CreateByQuestion"
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

	if err := s.tgUserService.UpdateState(ctx, user, constant.UserDefaultState); err != nil {
		return nil, err
	}

	return &sp, nil
}
