package githubclient

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v72/github"
)

// GitHub App installations have unique IDs, and each requires its own scoped authentication token.
// We can't use a single global GitHub client for multiple installations.
//
// To handle this, we cache per-installation GitHub clients.
// This avoids repeated token generation and ensures each client is properly scoped.
//
// Client initialization and caching is done in main.go to keep things centralized and efficient.

type ClientFactory struct {
	cfg   *config.Config
	cache map[int64]*github.Client
	mutex sync.Mutex
}

func NewClientFactory(cfg *config.Config) *ClientFactory {
	return &ClientFactory{
		cfg:   cfg,
		cache: make(map[int64]*github.Client),
	}
}

func (f *ClientFactory) GetClient(ctx context.Context, installationID int64) (*github.Client, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if client, exists := f.cache[installationID]; exists {
		return client, nil
	}

	tr, err := ghinstallation.New(http.DefaultTransport, f.cfg.GithubAppID, installationID, []byte(f.cfg.GithubAppPrivateKey))
	if err != nil {
		log.Printf("Error creating GitHub App transport: %v\n", err)
		return nil, err
	}
	client := github.NewClient(&http.Client{Transport: tr})
	f.cache[installationID] = client
	return client, nil
}
