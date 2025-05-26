package api

import (
	"net/http"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	ghclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/Mentro-Org/CodeLookout/internal/handlers"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(cfg *config.Config, ghClientFactory *ghclient.ClientFactory, aiClient llm.AIClient, dbPool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("ok")); err != nil {
				log.Printf("failed to write response: %v", err)
			}
		})
		service := handlers.NewWebhookHandlerService(cfg, ghClientFactory, aiClient, dbPool)
		r.Post("/webhook", service.HandleWebhook())
	})
	return r
}
