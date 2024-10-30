package steps

import "alertflow-runner/pkg/models"

func GetStepByActionName(steps []models.ExecutionSteps, actionName string) models.ExecutionSteps {
	for _, step := range steps {
		if step.ActionName == actionName {
			return step
		}
	}
	return models.ExecutionSteps{}
}
