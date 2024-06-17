package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/ten_tarot/internal/entity"
	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/go-playground/validator/v10"
	"github.com/sashabaranov/go-openai"
)

type Event struct {
	l      logger.Logger
	openai *openai.Client
	prompt string
}

func NewEvent(l logger.Logger, token, prompt string) *Event {
	return &Event{
		l:      l,
		openai: openai.NewClient(token),
		prompt: prompt,
	}
}

func (e *Event) CreateFromText(ctx context.Context, event *entity.Event, text string) error {
	const op = "./internal/service/event::CreateFromText"

	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: e.prompt + time.Now().Format("20060102T150405Z")},
		{Role: openai.ChatMessageRoleUser, Content: text},
	}
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: messages,
	}
	resp, err := e.openai.CreateChatCompletion(ctx, req)
	if err != nil {
		e.l.Error(fmt.Errorf("%s: %w", op, err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(resp.Choices) < 1 {
		e.l.Error(fmt.Errorf("%s: %w", op, err))
		return fmt.Errorf("%s: choices is empty", op)
	}

	choice := resp.Choices[0]
	jsonString := strings.Replace(strings.Replace(choice.Message.Content, "json", "", -1), "`", "", -1)

	if err = json.Unmarshal([]byte(jsonString), &event); err != nil {
		e.l.Error(fmt.Errorf("%s: %w", op, err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if err = validator.New().Struct(event); err != nil {
		e.l.Error(fmt.Errorf("%s: %w", op, err))
		return fmt.Errorf("%s: %w", op, failure.ErrValidation)
	}

	return nil
}
