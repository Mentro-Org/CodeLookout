package handlers

import (
	"log"
	"net/http"

	"github.com/Mentro-Org/CodeLookout/internal/core"
	"github.com/Mentro-Org/CodeLookout/internal/queue"

	"github.com/google/go-github/v72/github"
)

func HandleWebhook(appDeps *core.AppDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := github.ValidatePayload(r, []byte(appDeps.Config.WebhookSecret))
		if err != nil {
			http.Error(w, "Invalid payload signature", http.StatusUnauthorized)
			return
		}

		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			http.Error(w, "Could not parse webhook", http.StatusBadRequest)
			return
		}

		switch e := event.(type) {
		case *github.PullRequestEvent:
			if e.GetAction() == "opened" || e.GetAction() == "synchronize" {
				payload := queue.PRReviewTaskPayload{
					InstallationID: e.GetInstallation().GetID(),
					Owner:          e.GetRepo().GetOwner().GetLogin(),
					Repo:           e.GetRepo().GetName(),
					PRNumber:       e.GetNumber(),
					Title:          e.GetPullRequest().GetTitle(),
					Body:           e.GetPullRequest().GetBody(),
					CommitSHA:      e.GetPullRequest().GetHead().GetSHA(),
				}

				err := appDeps.TaskClient.EnqueueTask(payload)
				if err != nil {
					log.Printf("Error queuing task: %v", err)
					http.Error(w, "Queue error", http.StatusInternalServerError)
					return
				}
			}
		default:
			http.Error(w, "Unsupported event type", http.StatusNotImplemented)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		return
	}
}
