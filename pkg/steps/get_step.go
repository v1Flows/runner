package steps

import (
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

func GetStepByActionName(steps []bmodels.ExecutionSteps, actionName string) bmodels.ExecutionSteps {
	for _, step := range steps {
		if step.Action.Name == actionName {
			return step
		}
	}
	return bmodels.ExecutionSteps{}
}
