package api

import (
	"net/http"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	ghclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/Mentro-Org/CodeLookout/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg *config.Config, ghClientFactory *ghclient.ClientFactory, aiClient llm.AIClient,db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		service := handlers.NewWebhookHandlerService(cfg, ghClientFactory, aiClient, db)
		r.Post("/webhook", service.HandleWebhook())
	})

	return r
}
