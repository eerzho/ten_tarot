package service

import (
	"context"

	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
)

type SupportRequest struct {
	supportRequestRepo supportRequestRepo
	tgUserService      tgUserService
}

func NewSupportRequest(
	supportRequestRepo supportRequestRepo,
	tgUserService tgUserService,
) *SupportRequest {
	return &SupportRequest{
		supportRequestRepo: supportRequestRepo,
		tgUserService:      tgUserService,
	}
}

func (s *SupportRequest) CreateByUserQuestion(ctx context.Context, user *model.TGUser, question string) (*model.SupportRequest, error) {
	const op = "service.SupportRequest.CreateByQuestion"
	logger.Debug(
		op,
		logger.Any("user", user),
		logger.Any("question", question),
	)

	sp := model.SupportRequest{
		ChatID:   user.ChatID,
		Question: question,
	}

	if err := s.supportRequestRepo.Create(ctx, &sp); err != nil {
		return nil, err
	}

	if err := s.tgUserService.UpdateState(ctx, user, ""); err != nil {
		return nil, err
	}

	return &sp, nil
}
