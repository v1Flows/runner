package models

import ef_models "github.com/v1Flows/exFlow/services/backend/pkg/models"

type IncomingAlert struct {
	AlertData ef_models.Alerts `json:"alert"`
}

type IncomingAlerts struct {
	Alerts []ef_models.Alerts `json:"alerts"`
}
