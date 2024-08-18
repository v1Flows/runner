package models

import "encoding/json"

type IncomingPayload struct {
	PayloadData Payload `json:"payload"`
}

type Payload struct {
	Payload  json.RawMessage `json:"payload"`
	FlowID   string          `json:"flow_id"`
	RunnerID string          `json:"runner_id"`
}
