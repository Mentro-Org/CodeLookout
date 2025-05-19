package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mentro-Org/CodeLookout/internal/api"
	"github.com/Mentro-Org/CodeLookout/internal/config"
	githubclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/joho/godotenv"
)

func main() {
	// Only load .env in local/dev environments
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}

	cfg := config.Load()
	ghClientFactory := githubclient.NewClientFactory(cfg)
	aiClient, err := llm.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to get an AI Client %v", err)
	}

	// Setup router
	r := api.NewRouter(cfg, ghClientFactory, aiClient)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("Server is listening at http://localhost%s\n", addr)

	log.Fatal(http.ListenAndServe(addr, r))
}
