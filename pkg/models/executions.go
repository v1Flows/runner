package models

import (
	"time"

	"github.com/google/uuid"
)

type Execution struct {
	ID                  uuid.UUID `json:"id"`
	FlowID              string    `json:"flow_id"`
	PayloadID           string    `json:"payload_id"`
	RunnerID            string    `json:"runner_id"`
	Pending             bool      `json:"pending"`
	Running             bool      `json:"running"`
	Paused              bool      `json:"paused"`
	Canceled            bool      `json:"canceled"`
	NoPatternMatch      bool      `json:"no_pattern_match"`
	InteractionRequired bool      `json:"interaction_required"`
	Error               bool      `json:"error"`
	Finished            bool      `json:"finished"`
	CreatedAt           time.Time `json:"created_at"`
	ExecutedAt          time.Time `json:"executed_at"`
	FinishedAt          time.Time `json:"finished_at"`
}

type Executions struct {
	Executions []Execution `json:"executions"`
}
