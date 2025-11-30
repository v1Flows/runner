package models

import jf_models "github.com/JustLABv1/justflow/services/backend/pkg/models"

type IncomingAlert struct {
	AlertData jf_models.Alerts `json:"alert"`
}

type IncomingAlerts struct {
	Alerts []jf_models.Alerts `json:"alerts"`
}
