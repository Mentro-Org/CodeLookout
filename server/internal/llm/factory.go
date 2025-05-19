package llm

import (
	"fmt"

	"github.com/Mentro-Org/CodeLookout/internal/config"
)

func NewClient(cfg *config.Config) (AIClient, error) {
	switch cfg.AIProvider {
	case "openai":
		return NewOpenAIClient(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.AIProvider)
	}
}
