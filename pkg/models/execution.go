package models

import "github.com/JustLABv1/justflow/services/backend/pkg/models"

type IncomingExecution struct {
	ExecutionData models.Executions `json:"execution"`
}
