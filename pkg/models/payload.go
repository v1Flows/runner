package models

import bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

type IncomingPayload struct {
	PayloadData bmodels.Payloads `json:"payload"`
}
