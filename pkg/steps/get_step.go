package steps

import (
	"github.com/JustLABv1/justflow/services/backend/pkg/models"
)

func GetStepByActionName(steps []models.ExecutionSteps, actionName string) models.ExecutionSteps {
	for _, step := range steps {
		if step.Action.Name == actionName {
			return step
		}
	}
	return models.ExecutionSteps{}
}
