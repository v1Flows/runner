package internal_executions

import (
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/executions"
	shared_models "github.com/v1Flows/shared-library/pkg/models"
)

// SendFlowActionSteps sends all active flow actions to alertflow
func sendFlowActionSteps(cfg config.Config, execution shared_models.Executions, flow shared_models.Flows) (stepsWithIDs []shared_models.ExecutionSteps, err error) {
	for _, action := range flow.Actions {
		if !action.Active {
			continue
		}

		step := shared_models.ExecutionSteps{
			Action:      action,
			ExecutionID: execution.ID.String(),
			Status:      "pending",
		}

		// handle custom name
		if action.CustomName != "" {
			step.Action.Name = action.CustomName
		}

		stepID, err := executions.SendStep(cfg, execution, step)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		stepsWithIDs = append(stepsWithIDs, step)
	}

	return stepsWithIDs, nil
}
