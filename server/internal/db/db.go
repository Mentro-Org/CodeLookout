package db

import (
    "context"
    "log"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Mentro-Org/CodeLookout/internal/config" 
	
)

func ConnectDB(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
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
