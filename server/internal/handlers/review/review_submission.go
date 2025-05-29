package review

import (
	"log"

	"github.com/Mentro-Org/CodeLookout/internal/core"
	"github.com/google/go-github/v72/github"
)

type ReviewSubmission struct {
	Body     string
	Event    string // "APPROVE", "REQUEST_CHANGES", or "COMMENT"
	Comments []*github.DraftReviewComment
}

func (rs *ReviewSubmission) Execute(reviewCtx *core.ReviewContext) error {
	ctx := reviewCtx.Ctx
	client, err := reviewCtx.AppDeps.GHClientFactory.GetClient(ctx, reviewCtx.Payload.InstallationID)
	if err != nil {
		return err
	}

	review := &github.PullRequestReviewRequest{
		Body:     &rs.Body,
		Event:    &rs.Event,
		Comments: rs.Comments,
	}

	_, _, err = client.PullRequests.CreateReview(ctx,
		reviewCtx.Payload.Owner,
		reviewCtx.Payload.Repo,
		reviewCtx.Payload.PRNumber,
		review,
	)

	if err != nil {
		log.Printf("Failed to submit review: %v\n", err)
	}
	return err
}
