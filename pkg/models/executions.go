package models

import (
	"time"

	"github.com/google/uuid"
)

type Execution struct {
	ID                  uuid.UUID `json:"id"`
	FlowID              string    `json:"flow_id"`
	PayloadID           string    `json:"payload_id"`
	NoMatch             bool      `json:"no_match"`
	Running             bool      `json:"running"`
	Error               bool      `json:"error"`
	CreatedAt           time.Time `json:"created_at"`
	ExecutedAt          time.Time `json:"executed_at"`
	FinishedAt          time.Time `json:"finished_at"`
	RunnerID            string    `json:"runner_id"`
	Waiting             bool      `json:"waiting"`
	Paused              bool      `json:"paused"`
	Ghost               bool      `json:"ghost"`
	TotalSteps          int       `json:"total_steps"`
	Cancelled           bool      `json:"cancelled"`
	InteractionRequired bool      `json:"interaction_required"`
}

type Executions struct {
	Executions []Execution `json:"executions"`
}
