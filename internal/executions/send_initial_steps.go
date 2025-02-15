package internal_executions

import (
	"time"

	"github.com/AlertFlow/runner/pkg/executions"
	"github.com/google/uuid"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

// SendInitialSteps sends initial steps to alertflow
func SendInitialSteps(execution bmodels.Executions) (stepsWithIDs []bmodels.ExecutionSteps, err error) {
	var initialSteps = []bmodels.ExecutionSteps{
		{
			Action: bmodels.Actions{
				ID:          uuid.New(),
				Name:        "Runner Pick Up",
				Description: "Runner picked up the execution",
				Version:     "1.0.0",
				Icon:        "solar:rocket-2-bold-duotone",
				Type:        "runner_pick_up",
				Category:    "runner",
			},
			Messages: []string{
				execution.RunnerID + " picked up the execution",
			},
			Status:     "success",
			RunnerID:   execution.RunnerID,
			StartedAt:  time.Now(),
			FinishedAt: time.Now(),
		},
		{
			Action: bmodels.Actions{
				Plugin: "collect_data",
			},
			Status: "pending",
		},
		{
			Action: bmodels.Actions{
				Plugin: "pattern_check",
			},
			Status: "pending",
		},
		{
			Action: bmodels.Actions{
				Plugin: "actions_check",
			},
			Status: "pending",
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
