package models

import "github.com/v1Flows/exFlow/services/backend/pkg/models"

type IncomingExecutionSteps struct {
	StepsData []models.ExecutionSteps `json:"steps"`
}

type IncomingExecutionStep struct {
	StepData models.ExecutionSteps `json:"step"`
}
