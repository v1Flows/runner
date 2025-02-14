package models

import (
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

type Plugin struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Payload PayloadEndpoint
	Action  bmodels.Actions
}
