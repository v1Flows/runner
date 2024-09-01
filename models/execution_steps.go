package models

import (
	"time"

	"github.com/google/uuid"
)

type ExecutionSteps struct {
	ID             uuid.UUID `json:"id"`
	ExecutionID    string    `json:"execution_id"`
	ActionName     string    `json:"action_name"`
	ActionID       string    `json:"action_id"`
	ActionMessages []string  `json:"action_messages"`
	Error          bool      `json:"error"`
	Finished       bool      `json:"finished"`
	Paused         bool      `json:"paused"`
	StartedAt      time.Time `json:"started_at"`
	FinishedAt     time.Time `json:"finished_at"`
	NoResult       bool      `json:"no_result"`
	ParentID       string    `json:"parent_id"`
	IsHidden       bool      `json:"is_hidden"`
	NoPatternMatch bool      `json:"no_pattern_match"`
	Icon           string    `json:"icon"`
}
