package internal_executions

import (
	"github.com/AlertFlow/runner/pkg/executions"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

// SendFlowActionSteps sends all active flow actions to alertflow
func SendFlowActionSteps(execution bmodels.Executions, flow bmodels.Flows) (stepsWithIDs []bmodels.ExecutionSteps, err error) {
	for _, action := range flow.Actions {
		if !action.Active {
			continue
		}

		step := bmodels.ExecutionSteps{
			Action:      action,
			ExecutionID: execution.ID.String(),
			Status:      "pending",
		}

		// handle custom name
		if action.CustomName != "" {
			step.Action.Name = action.CustomName
		}

		stepID, err := executions.SendStep(execution, step)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		stepsWithIDs = append(stepsWithIDs, step)
	}

	return stepsWithIDs, nil
}
