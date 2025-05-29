package worker

import (
	"context"
	"log"

	"github.com/Mentro-Org/CodeLookout/internal/core"
	"github.com/Mentro-Org/CodeLookout/internal/handlers"
	"github.com/Mentro-Org/CodeLookout/internal/queue"
	"github.com/hibiken/asynq"
)

// This worker picks review jobs form queue and execute required actions
func RunWorker(ctx context.Context, appDeps *core.AppDeps) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: appDeps.Config.RedisAddress,
		},
		asynq.Config{
			Concurrency: appDeps.Config.WorkerConcurrency,
			Queues: map[string]int{
				queue.QueueName: 10,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TaskReviewPR, func(ctx context.Context, t *asynq.Task) error {
		return handlers.HandleReviewForPR(ctx, t, appDeps)
	})

	// Run in separate goroutine and stop when ctx is canceled
	// gracefully stop workers
	go func() {
		<-ctx.Done()
		log.Println("[Worker] shutting down Asynq server...")
		srv.Shutdown()
	}()

	log.Println("[Worker] Starting Asynq worker...")

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run worker server: %v", err)
	}

	log.Println("[Worker] Asynq worker stopped.")
}
