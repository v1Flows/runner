package models

import bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

type IncomingFlow struct {
	FlowData bmodels.Flows `json:"flow"`
}
