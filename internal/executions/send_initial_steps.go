package executions

import (
	"time"

	"gitlab.justlab.xyz/alertflow-public/runner/config"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"
)

var initialSteps = []models.ExecutionSteps{
	{
		ActionType:     "runner_pick_up",
		ActionName:     "Runner Pick Up",
		ActionMessages: []string{"Runner Picked Up Execution"},
		Icon:           "solar:rocket-2-bold-duotone",
		RunnerID:       config.Config.RunnerID,
		Pending:        false,
		Finished:       true,
		StartedAt:      time.Now(),
		FinishedAt:     time.Now(),
	},
	{
		ActionType:     "execution_start",
		ActionName:     "Execution Start",
		ActionMessages: []string{"Execution Started"},
		Icon:           "solar:rocket-2-bold-duotone",
		RunnerID:       config.Config.RunnerID,
		StartedAt:      time.Now(),
		Pending:        false,
		Finished:       true,
		FinishedAt:     time.Now(),
	},
	{
		ActionType: "collect_data",
		ActionName: "Collect Data",
		Icon:       "solar:inbox-archive-linear",
		Pending:    true,
	},
	{
		ActionType: "pattern_check",
		ActionName: "Pattern Check",
		Icon:       "solar:list-check-minimalistic-bold",
		Pending:    true,
	},
	{
		ActionType: "actions_check",
		ActionName: "Actions Check",
		Icon:       "solar:bolt-linear",
		Pending:    true,
	},
}

// SendInitialSteps sends initial steps to alertflow
func SendInitialSteps(execution models.Execution) (stepsWithIDs []models.ExecutionSteps, err error) {
	for i, step := range initialSteps {
		step.ExecutionID = execution.ID.String()
		stepID, err := SendStep(execution, step)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		initialSteps[i] = step
	}
	return initialSteps, nil
}
