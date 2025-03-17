package models

import (
	shared_models "github.com/v1Flows/shared-library/pkg/models"
)

type IncomingFlow struct {
	FlowData shared_models.Flows `json:"flow"`
}
