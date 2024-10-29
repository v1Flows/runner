package models

import (
	"time"

	"github.com/google/uuid"
)

type IncomingFlow struct {
	FlowData Flows `json:"flow"`
}

type Flows struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ProjectID          string    `json:"project_id"`
	RunnerID           string    `json:"runner_id"`
	ExecParallel       bool      `json:"exec_parallel"`
	Actions            []Actions `json:"actions"`
	Patterns           []Pattern `json:"patterns"`
	Maintenance        bool      `json:"maintenance"`
	MaintenanceMessage string    `json:"maintenance_message"`
	Disabled           bool      `json:"disabled"`
	DisabledReason     string    `json:"disabled_reason"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Actions struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Active      bool      `json:"active"`
	Params      []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"params"`
	CustomName        string `json:"custom_name"`
	CustomDescription string `json:"custom_description"`
}

type Pattern struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
