package models

import (
	"encoding/json"
	"time"
)

type Register struct {
	Registered                bool            `json:"registered"`
	AvailableActions          json.RawMessage `json:"available_actions"`
	AvailablePayloadInjectors json.RawMessage `json:"available_payload_injectors"`
	LastHeartbeat             time.Time       `json:"last_heartbeat"`
	RunnerVersion             string          `json:"runner_version"`
	Mode                      string          `json:"mode"`
}
