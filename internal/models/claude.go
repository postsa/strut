package models

import (
	"context"
	"errors"
	"fmt"
	anthropic "github.com/liushuangls/go-anthropic/v2"
	"os"
)

type Claude struct {
	client *anthropic.Client
	model  anthropic.Model // e.g., "claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307"
}

func NewClaude() (*Claude, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, errors.New("ANTHROPIC_API_KEY is required")
	}

	client := anthropic.NewClient(apiKey)

	return &Claude{
		client: client,
		model:  anthropic.ModelClaudeOpus4Dot0,
	}, nil
}

func (c *Claude) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     c.model,
		MaxTokens: 1024,
		Messages: []anthropic.Message{
			{
				Role:    anthropic.RoleUser,
				Content: []anthropic.MessageContent{{Text: &prompt, Type: anthropic.MessagesContentTypeText}},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate content with Claude: %w", err)
	}

	if len(resp.Content) == 0 {
		return "", errors.New("no content received from Claude API")
	}

	var generatedText string
	for _, contentBlock := range resp.Content {
		if contentBlock.Type == anthropic.MessagesContentTypeText {
			generatedText += *contentBlock.Text
		}
	}

	if generatedText == "" {
		return "", errors.New("no text content found in Claude API response")
	}

	return generatedText, nil
}

func (c *Claude) Close() {}
