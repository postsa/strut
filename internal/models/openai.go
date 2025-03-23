package models

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"os"
)

type OpenAi struct {
	client *openai.Client
}

func (o *OpenAi) GenerateContent(ctx context.Context, prompt string) (string, error) {
	content, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	resp := fmt.Sprintf("%v", content.Choices[0].Message.Content)
	return resp, nil
}

func NewOpenAi() (*OpenAi, error) {
	apiKey := os.Getenv("OPEN_AI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPEN_AI_API_KEY environment variable not set.")
		return nil, fmt.Errorf("OPEN_AI_API_KEY not set")
	}
	config := openai.DefaultConfig(apiKey)
	c := openai.NewClientWithConfig(config)
	return &OpenAi{
		client: c,
	}, nil
}

func (o *OpenAi) Close() {}
