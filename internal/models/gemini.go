package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Gemini encapsulates the Gemini API client.
type Gemini struct {
	genaiClient *genai.Client
	model       *genai.GenerativeModel
	chatSession *genai.ChatSession
}

// NewGemini creates a new Gemini API client.  Reads API key from environment.
func NewGemini(ctx context.Context) (*Gemini, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set.")
		return nil, fmt.Errorf("GEMINI API_KEY not set")
	}

	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := genaiClient.GenerativeModel("gemini-2.0-flash")

	cs := model.StartChat()

	return &Gemini{genaiClient: genaiClient, model: model, chatSession: cs}, nil
}

// GenerateContent sends a prompt to Gemini and returns the response.
func (c *Gemini) GenerateContent(ctx context.Context, prompt string) (string, error) {
	content, err := c.chatSession.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	resp := fmt.Sprintf("%v", content.Candidates[0].Content.Parts[0])

	return resp, nil
}

// Close closes the Gemini API client.
func (c *Gemini) Close() {
	c.genaiClient.Close()
}
