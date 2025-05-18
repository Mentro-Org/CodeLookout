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
	event := reviewCtx.Event
	ctx := reviewCtx.Ctx
	client, err := reviewCtx.GHClientFactory.GetClient(ctx, event.GetInstallation().GetID())
	if err != nil {
		return err
	}

	review := &github.PullRequestReviewRequest{
		Body:     &rs.Body,
		Event:    &rs.Event,
		Comments: rs.Comments,
	}

	_, _, err = client.PullRequests.CreateReview(ctx,
		event.GetRepo().GetOwner().GetLogin(),
		event.GetRepo().GetName(),
		event.GetNumber(),
		review,
	)

	if err != nil {
		log.Printf("Failed to submit review: %v\n", err)
	}
	return err
}
