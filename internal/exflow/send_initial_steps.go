package internal_exflow

import (
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/executions"
	"github.com/v1Flows/runner/pkg/platform"
	"github.com/v1Flows/shared-library/pkg/models"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

// SendInitialSteps sends initial steps to alertflow
func SendInitialSteps(cfg config.Config, actions []shared_models.Action, execution shared_models.Executions) (stepsWithIDs []shared_models.ExecutionSteps, err error) {
	var initialSteps = []shared_models.ExecutionSteps{
		{
			Action: shared_models.Action{
				Plugin: "collect_data",
				Params: []shared_models.Params{
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

	// get all current steps to modify the pickup step
	steps, err := executions.GetSteps(cfg, execution.ID.String(), targetPlatform)
	if err != nil {
		log.Error("Failed to get steps for execution: ", err)
		return
	}
	for _, step := range steps {
		if step.Action.Name == "Pick Up" {
			// modify the pickup step
			err = executions.UpdateStep(cfg, execution.ID.String(), models.ExecutionSteps{
				ID: step.ID,
				Messages: []models.Message{
					{
						Title: "Pick Up",
						Lines: []models.Line{
							{
								Content:   execution.RunnerID + " picked up the execution",
								Timestamp: time.Now(),
								Color:     "success",
							},
						},
					},
				},
				Status:     "success",
				RunnerID:   execution.RunnerID,
				FinishedAt: time.Now(),
			}, targetPlatform)
			if err != nil {
				return nil, err
			}
		}
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

		stepID, err := executions.SendStep(cfg, execution, step, targetPlatform)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		initialSteps[i] = step
	}
	return initialSteps, nil
}
