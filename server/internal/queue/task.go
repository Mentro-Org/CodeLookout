package queue

const TaskReviewPR = "review:pr"
const QueueName = "pr-review"

type PRReviewTaskPayload struct {
	InstallationID int64  `json:"installation_id"`
	Owner          string `json:"owner"`
	Repo           string `json:"repo"`
	PRNumber       int    `json:"pr_number"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CommitSHA      string `json:"commit_sha"`
}
