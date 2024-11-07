package executions

import (
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"
)

// SendFlowActionSteps sends all active flow actions to alertflow
func SendFlowActionSteps(execution models.Execution, flow models.Flows) (stepsWithIDs []models.ExecutionSteps, err error) {
	for _, action := range flow.Actions {
		step := models.ExecutionSteps{
			ActionID:    action.ID.String(),
			ActionType:  action.Type,
			ActionName:  action.Name,
			Icon:        action.Icon,
			ExecutionID: execution.ID.String(),
			Pending:     true,
		}

		// handle custom name
		if action.CustomName != "" {
			step.ActionName = action.CustomName
		}

		stepID, err := SendStep(execution, step)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		stepsWithIDs = append(stepsWithIDs, step)
	}

	return stepsWithIDs, nil
}
