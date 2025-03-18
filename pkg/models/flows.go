package models

import (
	af_models "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	ef_models "github.com/v1Flows/exFlow/services/backend/pkg/models"
	shared_models "github.com/v1Flows/shared-library/pkg/models"
)

type IncomingSharedFlow struct {
	FlowData shared_models.Flows `json:"flow"`
}

type IncomingAfFlow struct {
	FlowData af_models.Flows `json:"flow"`
}

type IncomingEfFlow struct {
	FlowData ef_models.Flows `json:"flow"`
}
