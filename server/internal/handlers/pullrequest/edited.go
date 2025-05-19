package pullrequest

import (
	"context"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	ghclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
    "github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/go-github/v72/github"
)

type PullRequestEditedHandler struct {
	Cfg             *config.Config
	AIClient        llm.AIClient
	GHClientFactory *ghclient.ClientFactory
	DB              *pgxpool.Pool
}

func (h *PullRequestEditedHandler) Handle(ctx context.Context, event *github.PullRequestEvent) error {
	return HandleReviewForPR(ctx, event, h.Cfg, h.GHClientFactory, h.AIClient, h.DB)
}
