package queue

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type TaskClient struct {
	client *asynq.Client
}

func NewTaskClient(redisAddr string) *TaskClient {
	return &TaskClient{
		client: asynq.NewClient(
			asynq.RedisClientOpt{
				Addr: redisAddr,
			},
		),
	}
}

func (tc *TaskClient) EnqueueTask(payload PRReviewTaskPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(TaskReviewPR, data)
	_, err = tc.client.Enqueue(
		task,
		asynq.Queue(QueueName),
		asynq.MaxRetry(5),
		asynq.Timeout(60*time.Second),
		asynq.Retention(24*time.Hour),
	)
	return err
}
