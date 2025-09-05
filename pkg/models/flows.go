package models

import (
	ef_models "github.com/v1Flows/exFlow/services/backend/pkg/models"
)

type IncomingFlow struct {
	FlowData ef_models.Flows `json:"flow"`
}
