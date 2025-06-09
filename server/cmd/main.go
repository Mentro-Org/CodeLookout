package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Mentro-Org/CodeLookout/internal/api"
	"github.com/Mentro-Org/CodeLookout/internal/config"
	"github.com/Mentro-Org/CodeLookout/internal/core"
	"github.com/Mentro-Org/CodeLookout/internal/db"
	githubclient "github.com/Mentro-Org/CodeLookout/internal/github"
	"github.com/Mentro-Org/CodeLookout/internal/llm"
	"github.com/Mentro-Org/CodeLookout/internal/queue"
	"github.com/Mentro-Org/CodeLookout/internal/worker"
	"github.com/joho/godotenv"
)

func initializeDependencies() (*core.AppDeps, error) {
	cfg := config.Load()
	ctx := context.Background()

	dbPool := db.ConnectDB(ctx, cfg)
	ghClientFactory := githubclient.NewClientFactory(cfg)

	aiClient, err := llm.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("AI client error: %w", err)
	}

	taskClient := queue.NewTaskClient(cfg.RedisAddress)

	return &core.AppDeps{
		Config:          cfg,
		GHClientFactory: ghClientFactory,
		AIClient:        aiClient,
		DBPool:          dbPool,
		TaskClient:      taskClient,
	}, nil
}

func startServer(ctx context.Context, deps *core.AppDeps) error {
	address := fmt.Sprintf(":%s", deps.Config.Port)
	srv := http.Server{
		Addr:    address,
		Handler: api.NewRouter(deps),
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down HTTP server...")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxTimeout); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("Server is listening at http://localhost%s\n", address)
	return srv.ListenAndServe()
}

func main() {
	// Only load .env in local/dev environments
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}

	appDeps, err := initializeDependencies()
	if err != nil {
		log.Fatalf("failed to initialize dependencies: %v", err)
	}
	defer appDeps.DBPool.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		// Whenever the program receives a SIGINT (Ctrl+C) or
		// SIGTERM (common in Docker/K8s), send it into channel c.
		// channel "c" subscribe to those signals.
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	// Start Asynq worker in background
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.RunWorker(ctx, appDeps)
	}()

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startServer(ctx, appDeps); err != nil {
			log.Printf("HTTP server exited: %v", err)
			cancel()
		}
	}()

	// Wait for everything to shut down
	wg.Wait()
	log.Println("App shutdown complete.")
}
