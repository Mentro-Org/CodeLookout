package pullrequest

import (
	"context"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	ghclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/google/go-github/v72/github"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestOpenedHandler struct {
	Cfg             *config.Config
	AIClient        llm.AIClient
	GHClientFactory *ghclient.ClientFactory
	DBPool          *pgxpool.Pool
}

func (h *PullRequestOpenedHandler) Handle(ctx context.Context, event *github.PullRequestEvent) error {
	return HandleReviewForPR(ctx, event, h.Cfg, h.GHClientFactory, h.AIClient, h.DBPool)
}
