package db

import (
	"context"
	"log"

	"github.com/Mentro-Org/CodeLookout/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set or is empty.")
	}
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatalf("Failed to ping DB: %v", err)
	}
	return pool
}
