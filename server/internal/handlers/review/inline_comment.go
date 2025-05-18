package review

import (
	"fmt"
	"log"

	"github.com/Mentro-Org/CodeLookout/internal/core"
	"github.com/google/go-github/v72/github"
)

type InlineComment struct {
	Body      string
	Path      string // file path
	StartLine int    // starting line in the diff
	Line      int    // ending line in the diff
}

func (ic *InlineComment) Execute(reviewCtx *core.ReviewContext) error {
	event := reviewCtx.Event
	ctx := reviewCtx.Ctx
	client, err := reviewCtx.GHClientFactory.GetClient(ctx, event.GetInstallation().GetID())
	if err != nil {
		return err
	}

	commitSHA := event.GetPullRequest().GetHead().GetSHA()

	comment := &github.PullRequestComment{
		Body:     github.Ptr(ic.Body),
		CommitID: github.Ptr(commitSHA),
		Path:     github.Ptr(ic.Path),
	}

	// single-line comment
	if ic.StartLine == ic.Line {
		comment.Line = github.Ptr(ic.Line)
	} else if ic.StartLine < ic.Line {
		// multi-line comment
		comment.StartLine = github.Ptr(ic.StartLine)
		comment.Line = github.Ptr(ic.Line)
	} else {
		return fmt.Errorf("start_line (%d) must be less than or equal to line (%d)", ic.StartLine, ic.Line)
	}

	_, _, err = client.PullRequests.CreateComment(ctx,
		event.GetRepo().GetOwner().GetLogin(),
		event.GetRepo().GetName(),
		event.GetNumber(),
		comment,
	)

	if err != nil {
		log.Printf("Failed to create inlmulti-line inline comment: %v\n", err)
	}
	return err
}
