package service

import (
	"context"
	"fmt"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/sashabaranov/go-openai"
)

type (
	Tarot interface {
		Oracle(ctx context.Context, question string, hand []model.Card) (string, error)
	}
	tarot struct {
		openai *openai.Client
		model  string
		prompt string
	}
)

func NewTarot(model, token, prompt string) Tarot {
	return &tarot{
		openai: openai.NewClient(token),
		model:  model,
		prompt: prompt,
	}
}

func (t *tarot) Oracle(ctx context.Context, question string, hand []model.Card) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: t.prompt},
		{
			Role: openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("My question is: %s\nI drew the following %d Tarot cards: %s\nCan you provide a detailed interpretation of each card and their meanings in this context?",
				question,
				len(hand),
				hand,
			),
		},
	}
	req := openai.ChatCompletionRequest{
		Model:    t.model,
		Messages: messages,
	}
	resp, err := t.openai.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) < 1 {
		return "", failure.ErrChoicesIsEmpty
	}

	return resp.Choices[0].Message.Content, nil
}
