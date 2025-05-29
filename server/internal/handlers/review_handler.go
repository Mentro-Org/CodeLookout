package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Mentro-Org/CodeLookout/internal/core"
	"github.com/Mentro-Org/CodeLookout/internal/handlers/review"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/Mentro-Org/CodeLookout/internal/queue"
	"github.com/hibiken/asynq"
)

func HandleReviewForPR(ctx context.Context, t *asynq.Task, appDeps *core.AppDeps) error {
	fmt.Println("received job from Queue")

	var payload queue.PRReviewTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Printf("Received a pull request event for #%d\n", payload.PRNumber)
	// Get GitHub client
	client, err := appDeps.GHClientFactory.GetClient(ctx, payload.InstallationID)
	if err != nil {
		return err
	}

	// Get changed files
	files, _, err := client.PullRequests.ListFiles(ctx,
		payload.Owner,
		payload.Repo,
		payload.PRNumber,
		nil,
	)
	if err != nil {
		return err
	}

	// Build prompt using file diffs
	promptText := llm.BuildPRReviewPrompt(&payload, files)

	var response string
	if appDeps.Config.AppEnv == "development" {
		// comment this out when pushing code
		response, err = appDeps.AIClient.GenerateSampleReviewForPR()
		if err != nil {
			return err
		}
	} else {
		// Call OpenAI with this prompt
		response, err = appDeps.AIClient.GenerateReviewForPR(ctx, promptText)
		if err != nil {
			return err
		}
	}

	err = review.HandleReviewResponseFromAI(ctx, payload, appDeps, response)
	if err != nil {
		return err
	}

	return nil
}
