package models

import (
	"time"

	ef_models "github.com/v1Flows/exFlow/services/backend/pkg/models"
)

type Register struct {
	ID             string                     `json:"id"`
	Registered     bool                       `json:"registered"`
	LastHeartbeat  time.Time                  `json:"last_heartbeat"`
	Version        string                     `json:"version"`
	Mode           string                     `json:"mode"`
	Plugins        []ef_models.Plugin         `json:"plugins"`
	Actions        []ef_models.Action         `json:"actions"`
	AlertEndpoints []ef_models.AlertEndpoints `json:"alert_endpoints"`
}
