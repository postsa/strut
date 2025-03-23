package models

import "context"

type ChatClient interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
	Close()
}
