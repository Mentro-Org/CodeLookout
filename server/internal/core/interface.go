package core

import (
	"context"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	ghclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/Mentro-Org/CodeLookout/internal/queue"
	"github.com/google/go-github/v72/github"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppDeps struct {
	Config          *config.Config
	GHClientFactory *ghclient.ClientFactory
	AIClient        llm.AIClient
	DBPool          *pgxpool.Pool
	TaskClient      *queue.TaskClient
}

type PullRequestHandler interface {
	Handle(ctx context.Context, event *github.PullRequestEvent) error
}

type ReviewAction interface {
	Execute(reviewCtx *ReviewContext) error
}

type ReviewContext struct {
	Ctx     context.Context
	Payload queue.PRReviewTaskPayload
	AppDeps *AppDeps
}
