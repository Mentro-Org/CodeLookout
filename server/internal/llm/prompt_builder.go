package llm

import (
	"fmt"

	"github.com/google/go-github/v72/github"
)

func BuildPRReviewPrompt(event *github.PullRequestEvent, files []*github.CommitFile) string {
	title := event.GetPullRequest().GetTitle()
	body := event.GetPullRequest().GetBody()

	prompt := fmt.Sprintf(`You are a senior software engineer performing a code review.

## PR Metadata
- **Title**: %s
- **Description**: %s

## Objective
Provide structured review feedback based on the following file diffs. The feedback will be used to generate review comments in GitHub.

## Output JSON Format:
{
  "summary": "High-level review summary (optional)",
  "action": "COMMENT | APPROVE | REQUEST_CHANGES",
  "comments": [
    {
      "path": "file/path.go",
      "comments": [
        {
          "line": { "s": 12, "e": 14 },
          "body": "Feedback for this line range",
          "category": ["security", "performance"]
        },
        {
          "line": { "s": 30, "e": 30 },
          "body": "Single-line comment",
          "category": ["style"]
        }
      ]
    }
  ]
}

- Use { "s": x, "e": y } to represent line ranges. If it's one line, make s == e.
- Categories can include: "security", "performance", "bug", "style", "readability", "lint", etc.
- Be concise, helpful, and professional.

`, title, body)

	for _, f := range files {
		prompt += fmt.Sprintf(`

---
**File**: %s
**Patch**:
%s`, f.GetFilename(), f.GetPatch())
	}

	prompt += `

---
Now analyze the code changes and respond with the JSON review object as specified.`

	return prompt
}
