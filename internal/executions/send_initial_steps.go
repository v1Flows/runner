package internal_executions

import (
	"time"

	"github.com/AlertFlow/runner/pkg/executions"
	"github.com/AlertFlow/runner/pkg/models"
)

// SendInitialSteps sends initial steps to alertflow
func SendInitialSteps(execution models.Execution) (stepsWithIDs []models.ExecutionSteps, err error) {
	var initialSteps = []models.ExecutionSteps{
		{
			ActionType: "runner_pick_up",
			ActionName: "Runner Pick Up",
			ActionMessages: []string{
				execution.RunnerID + " picked up the execution",
			},
			Icon:       "solar:rocket-2-bold-duotone",
			RunnerID:   execution.RunnerID,
			Pending:    false,
			Finished:   true,
			StartedAt:  time.Now(),
			FinishedAt: time.Now(),
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

	for i, step := range initialSteps {
		step.ExecutionID = execution.ID.String()

		stepID, err := executions.SendStep(execution, step)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		initialSteps[i] = step
	}
	return initialSteps, nil
}
