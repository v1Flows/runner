package models

import (
	"time"

	"github.com/google/uuid"
)

type IncomingFlow struct {
	FlowData   Flows         `json:"flow"`
	ActionData []FlowActions `json:"actions"`
}

type Flows struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	ProjectID           string    `json:"project_id"`
	RunnerID            string    `json:"runner_id"`
	Disabled            bool      `json:"disabled"`
	DisabledReason      string    `json:"disabled_reason"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	MaintenanceRequired bool      `json:"maintenance_required"`
	MaintenanceMessage  string    `json:"maintenance_message"`
}

type FlowActions struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	FlowID          string    `json:"flow_id"`
	Status          bool      `json:"status"`
	Actions         []string  `json:"actions"`
	ExecParallel    bool      `json:"exec_parallel"`
	MatchPatterns   []Pattern `json:"match_patterns"`
	ExcludePatterns []Pattern `json:"exclude_patterns"`
	CreatedAt       string    `json:"created_at"`
}

type Pattern struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
