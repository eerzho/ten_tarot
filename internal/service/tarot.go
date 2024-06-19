package service

import (
	"context"
	"fmt"

	"github.com/eerzho/ten_tarot/internal/entity"
	"github.com/sashabaranov/go-openai"
)

type Tarot struct {
	openai *openai.Client
	prompt string
}

func NewTarot(token, prompt string) *Tarot {
	return &Tarot{
		openai: openai.NewClient(token),
		prompt: prompt,
	}
}

func (t *Tarot) Oracle(ctx context.Context, question string, hand []entity.Card) (string, error) {
	const op = "./internal/service/tarot::Oracle"

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
		Model:    openai.GPT4o,
		Messages: messages,
	}
	resp, err := t.openai.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if len(resp.Choices) < 1 {
		return "", fmt.Errorf("%s: choices is empty", op)
	}

	choice := resp.Choices[0]

	return choice.Message.Content, nil
}
