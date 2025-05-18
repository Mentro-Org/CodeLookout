package llm

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

type AIClient interface {
	GenerateReviewForPR(ctx context.Context, prompt string) (string, error)
	GenerateSampleReviewForPR() (string, error)
}

type OpenAIClient struct {
	client *openai.Client
	model  string
}

var (
	instance *OpenAIClient
	once     sync.Once
)

// NewOpenAIClient initializes the singleton instance
func NewOpenAIClient(cfg *config.Config) *OpenAIClient {
	once.Do(func() {
		c := openai.NewClient(cfg.OpenAIKey)
		instance = &OpenAIClient{client: c, model: openai.GPT4}
	})
	return instance
}

func (c *OpenAIClient) GenerateReviewForPR(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are a senior code reviewer. Be concise."},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// use this for development(no need to call AI API to get review json)
func (c *OpenAIClient) GenerateSampleReviewForPR() (string, error) {
	rootDir, _ := os.Getwd()
	jsonPath := filepath.Join(rootDir, "data", "openai-review.json")
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	return string(bytes), nil
}
