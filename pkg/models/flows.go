package models

import (
	jf_models "github.com/JustLABv1/justflow/services/backend/pkg/models"
)

type IncomingFlow struct {
	FlowData jf_models.Flows `json:"flow"`
}
