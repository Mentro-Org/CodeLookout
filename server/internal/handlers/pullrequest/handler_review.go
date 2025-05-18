package pullrequest

import (
	"context"
	"log"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	ghclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/Mentro-Org/CodeLookout/internal/handlers/review"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/google/go-github/v72/github"
)

func HandleReviewForPR(ctx context.Context, event *github.PullRequestEvent, cfg *config.Config, ghClientFactory *ghclient.ClientFactory, aIClient llm.AIClient) error {
	log.Printf("Received a pull request event for #%d\n", event.GetNumber())
	// Get GitHub client
	client, err := ghClientFactory.GetClient(ctx, event.GetInstallation().GetID())
	if err != nil {
		return err
	}

	// Get changed files
	files, _, err := client.PullRequests.ListFiles(ctx,
		event.GetRepo().GetOwner().GetLogin(),
		event.GetRepo().GetName(),
		event.GetNumber(),
		nil,
	)
	if err != nil {
		return err
	}

	// Build prompt using file diffs
	promptText := llm.BuildPRReviewPrompt(event, files)

	var response string
	if cfg.AppEnv == "development" {
		// comment this out when pushing code
		response, err = aIClient.GenerateSampleReviewForPR()
		if err != nil {
			return err
		}
	} else {
		// Call OpenAI with this prompt
		response, err = aIClient.GenerateReviewForPR(ctx, promptText)
		if err != nil {
			return err
		}
	}

	err = review.HandleReviewResponseFromAI(ctx, event, cfg, ghClientFactory, response)
	if err != nil {
		return err
	}

	return nil
}
