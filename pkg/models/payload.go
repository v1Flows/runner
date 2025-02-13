package models

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type IncomingPayload struct {
	PayloadData Payload `json:"payload"`
}

type Payload struct {
	Payload  json.RawMessage `json:"payload"`
	FlowID   string          `json:"flow_id"`
	RunnerID string          `json:"runner_id"`
	Endpoint string          `json:"endpoint"`
}

type PayloadEndpoint struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
	Version  string `json:"version"`
}

type PayloadHandler interface {
	Init() PayloadEndpoint
	Handle(context *gin.Context)
}
