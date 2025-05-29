package internal_alertflow

import (
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/executions"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

// SendInitialSteps sends initial steps to alertflow
func SendInitialSteps(cfg *config.Config, actions []shared_models.Action, execution shared_models.Executions, alertID string) (stepsWithIDs []shared_models.ExecutionSteps, err error) {
	var initialSteps = []shared_models.ExecutionSteps{
		{
			Action: shared_models.Action{
				Name:        "Runner Pick Up",
				Description: "Runner picked up the execution",
				Version:     "1.0.0",
				Icon:        "solar:rocket-2-bold-duotone",
				Category:    "runner",
			},
			Messages: []shared_models.Message{
				{
					Title: "Runner Pick Up",
					Lines: []shared_models.Line{
						{
							Content: execution.RunnerID + " picked up the execution",
						},
					},
				},
			},
			Status:     "success",
			RunnerID:   execution.RunnerID,
			CreatedAt:  time.Now(),
			StartedAt:  time.Now(),
			FinishedAt: time.Now(),
		},
		{
			Action: shared_models.Action{
				Plugin: "collect_data",
				Params: []shared_models.Params{
					{
						Key:   "AlertID",
						Value: alertID,
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
			Action: shared_models.Action{
				Plugin: "pattern_check",
			},
			Status:    "pending",
			CreatedAt: time.Now(),
		},
		{
			Action: shared_models.Action{
				Plugin: "actions_check",
			},
			Status:    "pending",
			CreatedAt: time.Now(),
		},
	}

	targetPlatform, ok := platform.GetPlatformForExecution(execution.ID.String())
	if !ok {
		log.Error("Failed to get platform")
		return
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

		stepID, err := executions.SendStep(nil, execution, step, targetPlatform)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		initialSteps[i] = step
	}
	return initialSteps, nil
}
