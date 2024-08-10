package models

import (
	"time"

	"github.com/google/uuid"
)

type Flows struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	ProjectID           string    `json:"project_id"`
	RunnerID            string    `json:"runner_id"`
	Disabled            bool      `json:"disabled"`
	DisabledReason      string    `json:"disabled_reason"`
	Actions             []Action  `json:"actions"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	MaintenanceRequired bool      `json:"maintenance_required"`
	MaintenanceMessage  string    `json:"maintenance_message"`
}

type Action struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	Action         string    `json:"action"`
	MatchPattern   []Pattern `json:"match_patterns"`
	ExcludePattern []Pattern `json:"exclude_patterns"`
}

type Pattern struct {
	Key     string `json:"key"`
	Group   string `json:"group"`
	Value   string `json:"value"`
	ReactOn string `json:"react_on"`
}
