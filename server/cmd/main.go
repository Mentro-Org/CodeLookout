package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mentro-Org/CodeLookout/internal/api"
	"github.com/Mentro-Org/CodeLookout/internal/config"
	"github.com/Mentro-Org/CodeLookout/internal/db"
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

	ctx := context.Background()
	// Connect to the database
	// Initialize the database connection
	dbPool := db.ConnectDB(ctx, cfg)
	defer dbPool.Close()
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database")

	ghClientFactory := githubclient.NewClientFactory(cfg)
	aiClient := llm.NewOpenAIClient(cfg)

	// Setup router
	r := api.NewRouter(cfg, ghClientFactory, aiClient, dbPool)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("Server is listening at http://localhost%s\n", addr)

	log.Fatal(http.ListenAndServe(addr, r))
}
