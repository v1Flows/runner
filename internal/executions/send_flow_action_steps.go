package internal_executions

import (
	"github.com/JustLABv1/justflow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/executions"
)

// SendFlowActionSteps sends all active flow actions to alertflow
func sendFlowActionSteps(cfg *config.Config, execution models.Executions, flow models.Flows) (stepsWithIDs []models.ExecutionSteps, err error) {
	for _, action := range flow.Actions {
		if !action.Active {
			continue
		}

		step := models.ExecutionSteps{
			Action:      action,
			ExecutionID: execution.ID.String(),
			Status:      "pending",
		}

		// handle custom name
		if action.CustomName != "" {
			step.Action.Name = action.CustomName
		}

		stepID, err := executions.SendStep(nil, execution, step)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		stepsWithIDs = append(stepsWithIDs, step)
	}

	return stepsWithIDs, nil
}
