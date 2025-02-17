package internal_executions

import (
	"time"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/executions"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

// SendInitialSteps sends initial steps to alertflow
func SendInitialSteps(cfg config.Config, actions []models.Actions, execution models.Executions) (stepsWithIDs []models.ExecutionSteps, err error) {
	var initialSteps = []models.ExecutionSteps{
		{
			Action: models.Actions{
				Name:        "Runner Pick Up",
				Description: "Runner picked up the execution",
				Version:     "1.0.0",
				Icon:        "solar:rocket-2-bold-duotone",
				Category:    "runner",
			},
			Messages: []string{
				execution.RunnerID + " picked up the execution",
			},
			Status:     "success",
			RunnerID:   execution.RunnerID,
			CreatedAt:  time.Now(),
			StartedAt:  time.Now(),
			FinishedAt: time.Now(),
		},
		{
			Action: models.Actions{
				Plugin: "collect_data",
				Params: []models.Params{
					{
						Key:   "PayloadID",
						Value: execution.PayloadID,
					},
					{
						Key:   "FlowID",
						Value: execution.FlowID,
					},
					{
						Key:   "LogData",
						Value: "false",
					},
				},
			},
			Status:    "pending",
			CreatedAt: time.Now(),
		},
		{
			Action: models.Actions{
				Plugin: "pattern_check",
			},
			Status:    "pending",
			CreatedAt: time.Now(),
		},
		{
			Action: models.Actions{
				Plugin: "actions_check",
			},
			Status:    "pending",
			CreatedAt: time.Now(),
		},
	}

	for i, step := range initialSteps {
		step.ExecutionID = execution.ID.String()

		// get action plugin info
		for _, action := range actions {
			if action.Plugin == step.Action.Plugin {
				if step.Action.Name == "" || step.Action.Description == "" {
					step.Action.Name = action.Name
					step.Action.Description = action.Description
					step.Action.Version = action.Version
					step.Action.Icon = action.Icon
					step.Action.Category = action.Category
				}
			}
		}

		stepID, err := executions.SendStep(cfg, execution, step)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		initialSteps[i] = step
	}
	return initialSteps, nil
}
