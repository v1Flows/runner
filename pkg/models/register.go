package models

import (
	"time"

	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

type Register struct {
	ID               string                     `json:"id"`
	Registered       bool                       `json:"registered"`
	LastHeartbeat    time.Time                  `json:"last_heartbeat"`
	Version          string                     `json:"version"`
	Mode             string                     `json:"mode"`
	Plugins          []bmodels.Plugins          `json:"plugins"`
	Actions          []bmodels.Actions          `json:"actions"`
	PayloadEndpoints []bmodels.PayloadEndpoints `json:"endpoints"`
}
