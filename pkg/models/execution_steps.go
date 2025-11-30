package models

import "github.com/JustLABv1/justflow/services/backend/pkg/models"

type IncomingExecutionSteps struct {
	StepsData []models.ExecutionSteps `json:"steps"`
}

type IncomingExecutionStep struct {
	StepData models.ExecutionSteps `json:"step"`
}
