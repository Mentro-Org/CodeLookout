package review

import (
	"log"

	"github.com/Mentro-Org/CodeLookout/internal/core"
	"github.com/google/go-github/v72/github"
)

type GeneralComment struct {
	Message string
}

func (gc *GeneralComment) Execute(reviewCtx *core.ReviewContext) error {
	ctx := reviewCtx.Ctx
	log.Printf("Received a pull request event for #%d\n", reviewCtx.Payload.PRNumber)

	client, err := reviewCtx.AppDeps.GHClientFactory.GetClient(ctx, reviewCtx.Payload.InstallationID)
	if err != nil {
		return err
	}

	_, _, err = client.Issues.CreateComment(
		ctx,
		reviewCtx.Payload.Owner,
		reviewCtx.Payload.Repo,
		reviewCtx.Payload.PRNumber,
		&github.IssueComment{Body: &gc.Message},
	)

	if err != nil {
		log.Printf("Failed to create general comment: %v\n", err)
	}
	return err
}
