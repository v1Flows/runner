package models

import (
	"encoding/json"
	"time"
)

type Register struct {
	ID               string          `json:"id"`
	Registered       bool            `json:"registered"`
	LastHeartbeat    time.Time       `json:"last_heartbeat"`
	Version          string          `json:"version"`
	Mode             string          `json:"mode"`
	Plugins          json.RawMessage `json:"plugins"`
	Actions          json.RawMessage `json:"actions"`
	PayloadEndpoints json.RawMessage `json:"payload_endpoints"`
}
