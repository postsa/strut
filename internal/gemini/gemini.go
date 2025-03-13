package gemini

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Client encapsulates the Gemini API client.
type Client struct {
	genaiClient *genai.Client
	model       *genai.GenerativeModel
}

// NewClient creates a new Gemini API client.  Reads API key from environment.
func NewClient(ctx context.Context) (*Client, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable not set.")
		return nil, fmt.Errorf("API_KEY not set")
	}

	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := genaiClient.GenerativeModel("gemini-2.0-flash") // You can make this configurable

	return &Client{genaiClient: genaiClient, model: model}, nil
}

// GenerateContent sends a prompt to Gemini and returns the response.
func (c *Client) GenerateContent(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Close closes the Gemini API client.
func (c *Client) Close() {
	c.genaiClient.Close()
}
