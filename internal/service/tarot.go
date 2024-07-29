package service

import (
	"context"
	"fmt"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/sashabaranov/go-openai"
)

type (
	Tarot struct {
		openaiClient *openai.Client
		modelName    string
		systemPrompt string
	}
)

func NewTarot(modelName, apiToken, systemPrompt string) *Tarot {
	return &Tarot{
		openaiClient: openai.NewClient(apiToken),
		modelName:    modelName,
		systemPrompt: systemPrompt,
	}
}

func (ts *Tarot) Oracle(ctx context.Context, userQuestion string, drawnCards []model.Card) (string, error) {
	const op = "service.Tarot.Oracle"
	logger.Debug(
		op,
		logger.Any("userQuestion", userQuestion),
		logger.Any("drawnCardsCount", len(drawnCards)),
	)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: ts.systemPrompt,
		},
		{
			Role: openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(
				"My question is: %s\nI drew the following %d Tarot cards: %s\nCan you provide a detailed interpretation of each card and their meanings in this context?",
				userQuestion,
				len(drawnCards),
				drawnCards,
			),
		},
	}
	req := openai.ChatCompletionRequest{
		Model:    ts.modelName,
		Messages: messages,
	}
	resp, err := ts.openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) < 1 {
		return "", failure.ErrChoicesIsEmpty
	}

	return resp.Choices[0].Message.Content, nil
}
