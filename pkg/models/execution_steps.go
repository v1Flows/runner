package models

import bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

type IncomingExecutionSteps struct {
	StepsData []bmodels.ExecutionSteps `json:"steps"`
}

type IncomingExecutionStep struct {
	StepData bmodels.ExecutionSteps `json:"step"`
}
