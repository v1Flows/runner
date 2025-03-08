package models

import bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

type IncomingAlert struct {
	AlertData bmodels.Alerts `json:"alert"`
}

type IncomingAlerts struct {
	Alerts []bmodels.Alerts `json:"alerts"`
}
