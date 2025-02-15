package models

import bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

type IncomingPayload struct {
	PayloadData bmodels.Payloads `json:"payload"`
}

type PayloadEndpoint struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
	Version  string `json:"version"`
}
