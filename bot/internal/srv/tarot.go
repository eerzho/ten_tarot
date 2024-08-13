package srv

import (
	"bot/internal/def"
	"bot/internal/dto"
	"context"
	"fmt"
	"log/slog"

	"github.com/sashabaranov/go-openai"
)

type Tarot struct {
	lg           *slog.Logger
	openAIClient *openai.Client
	modelName    string
	systemPrompt string
}

func NewTarot(lg *slog.Logger, modelName, apiToken, systemPrompt string) *Tarot {
	return &Tarot{
		lg:           lg,
		openAIClient: openai.NewClient(apiToken),
		modelName:    modelName,
		systemPrompt: systemPrompt,
	}
}

func (t *Tarot) Oracle(ctx context.Context, userQuestion string, drawnCards []dto.Card) (string, error) {
	const op = "srv.Tarot.Oracle"
	t.lg.Debug(
		op,
		slog.String("userQuestion", userQuestion),
		slog.Any("drawnCardsCount", len(drawnCards)),
	)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: t.systemPrompt,
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
		Model:    t.modelName,
		Messages: messages,
	}
	resp, err := t.openAIClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) < 1 {
		return "", def.ErrChoicesIsEmpty
	}

	return resp.Choices[0].Message.Content, nil
}
