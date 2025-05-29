package review

import (
	"context"
	"log"

	"github.com/Mentro-Org/CodeLookout/internal/core"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/Mentro-Org/CodeLookout/internal/queue"
)

func HandleReviewResponseFromAI(ctx context.Context, payload queue.PRReviewTaskPayload, appDeps *core.AppDeps, aiJsonResponse string) error {
	reviewResp, err := llm.ParseReviewResponse(aiJsonResponse)
	if err != nil {
		return err
	}

	reviewCtx := core.ReviewContext{
		Ctx:     ctx,
		Payload: payload,
		AppDeps: appDeps,
	}

	// Inline comments (multi-line support)
	for _, fileGroup := range reviewResp.Comments {
		for _, c := range fileGroup.Comments {
			// Format the comment body to include the line range if > 1 line
			commentBody := c.Body
			inline := &InlineComment{
				Body:      commentBody,
				Path:      fileGroup.Path,
				StartLine: c.Line.S,
				Line:      c.Line.E,
			}
			if err := inline.Execute(&reviewCtx); err != nil {
				log.Printf("Error posting inline comment on %s:%d: %v", fileGroup.Path, c.Line.S, err)
			}
		}
	}

	// Submit overall review event (APPROVE, REQUEST_CHANGES, COMMENT)
	reviewSubmission := &ReviewSubmission{
		Body:  reviewResp.Summary,
		Event: reviewResp.Action,
	}

	if err := reviewSubmission.Execute(&reviewCtx); err != nil {
		log.Printf("Error submitting review: %v", err)
		return err
	}

	return nil
}
