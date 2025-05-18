package core

import (
	"context"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	ghclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/google/go-github/v72/github"
)

type PullRequestHandler interface {
	Handle(ctx context.Context, event *github.PullRequestEvent) error
}

type ReviewAction interface {
	Execute(reviewCtx *ReviewContext) error
}

type ReviewContext struct {
	Ctx             context.Context
	Event           *github.PullRequestEvent
	Cfg             *config.Config
	GHClientFactory *ghclient.ClientFactory
}
