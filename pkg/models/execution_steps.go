package models

import (
	shared_models "github.com/v1Flows/shared-library/pkg/models"
)

type IncomingExecutionSteps struct {
	StepsData []shared_models.ExecutionSteps `json:"steps"`
}

type IncomingExecutionStep struct {
	StepData shared_models.ExecutionSteps `json:"step"`
}
