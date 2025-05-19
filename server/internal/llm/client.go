package llm

import "context"

type AIClient interface {
	GenerateReviewForPR(ctx context.Context, prompt string) (string, error)
	GenerateSampleReviewForPR() (string, error)
}
