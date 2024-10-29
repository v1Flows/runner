package models

import (
	"time"

	"github.com/google/uuid"
)

type IncomingExecutionSteps struct {
	StepsData []ExecutionSteps `json:"steps"`
}

type IncomingExecutionStep struct {
	StepData ExecutionSteps `json:"step"`
}

type ExecutionSteps struct {
	ID                  uuid.UUID `json:"id"`
	ExecutionID         string    `json:"execution_id"`
	ActionID            string    `json:"action_id"`
	ActionName          string    `json:"action_name"`
	ActionType          string    `json:"action_type"`
	ActionMessages      []string  `json:"action_messages"`
	RunnerID            string    `json:"runner_id"`
	ParentID            string    `json:"parent_id"`
	Icon                string    `json:"icon"`
	Interactive         bool      `json:"interactive"`
	Interacted          bool      `json:"interacted"`
	InteractionApproved bool      `json:"interaction_approved"`
	InteractionRejected bool      `json:"interaction_rejected"`
	InteractedBy        string    `json:"interacted_by"`
	InteractedAt        time.Time `json:"interacted_at"`
	IsHidden            bool      `json:"is_hidden"`
	Pending             bool      `json:"pending"`
	Running             bool      `json:"running"`
	Paused              bool      `json:"paused"`
	Canceled            bool      `json:"canceled"`
	CanceledBy          string    `json:"canceled_by"`
	CanceledAt          time.Time `json:"canceled_at"`
	NoPatternMatch      bool      `json:"no_pattern_match"`
	NoResult            bool      `json:"no_result"`
	Skipped             bool      `json:"skipped"`
	Error               bool      `json:"error"`
	Finished            bool      `json:"finished"`
	CreatedAt           time.Time `json:"created_at"`
	StartedAt           time.Time `json:"started_at"`
	FinishedAt          time.Time `json:"finished_at"`
}
