package llm

import (
	"encoding/json"
)

type PRReviewResponse struct {
	Summary  string      `json:"summary,omitempty"`
	Action   string      `json:"action"` // COMMENT, APPROVE, REQUEST_CHANGES
	Comments []FileGroup `json:"comments"`
}

type FileGroup struct {
	Path     string        `json:"path"`
	Comments []LineComment `json:"comments"`
}

type LineComment struct {
	Line     LineRange `json:"line"`
	Body     string    `json:"body"`
	Category []string  `json:"category,omitempty"`
}

type LineRange struct {
	S int `json:"s"` // start line
	E int `json:"e"` // end line
}

func ParseReviewResponse(jsonStr string) (*PRReviewResponse, error) {
	var response PRReviewResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
