package processing

import (
	"alertflow-runner/functions/actions"
	"alertflow-runner/functions/executions"
	"alertflow-runner/functions/flow"
	"alertflow-runner/functions/payload"
	"alertflow-runner/handlers/config"
	"alertflow-runner/handlers/variables"
	"alertflow-runner/models"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartProcessing(execution models.Execution) {
	execution.RunnerID = config.Config.RunnerID
	execution.Waiting = false
	execution.Running = true
	execution.ExecutedAt = time.Now()
	execution.TotalSteps = 2

	err := executions.Update(execution)
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// set runner picked up step
	_, err = executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Runner Pick Up",
		ActionMessages: []string{"Waiting for Runner to pick up Execution", "Runner picked up execution"},
		StartedAt:      execution.CreatedAt,
		Finished:       true,
		FinishedAt:     time.Now(),
		Icon:           "solar:rocket-2-bold-duotone",
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// collect data step
	collectDataStep, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionMessages: []string{"Collecting Data"},
		ActionName:     "Collect Data",
		StartedAt:      time.Now(),
		Icon:           "solar:inbox-archive-linear",
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// get flow data
	collectFlowDataStep, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Get Flow Data",
		ActionMessages: []string{"Requesting Flow Data from API"},
		StartedAt:      time.Now(),
		ParentID:       collectDataStep.ID.String(),
		IsHidden:       true,
		Icon:           "solar:book-bookmark-broken",
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	flowData, flowDataErr := flow.GetFlowData(execution)

	if flowDataErr != nil {
		err := executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             collectFlowDataStep.ID,
			ActionMessages: []string{"Failed to get Flow Data"},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			executions.EndWithError(execution)
			return
		}

		executions.EndWithError(execution)
		return
	}

	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             collectFlowDataStep.ID,
		ActionMessages: []string{"Flow Data received"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// get payload data
	collectPayloadDataStep, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Get Payload Data",
		ActionMessages: []string{"Requesting Payload Data from API"},
		StartedAt:      time.Now(),
		ParentID:       collectDataStep.ID.String(),
		IsHidden:       true,
		Icon:           "solar:letter-opened-broken",
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	payloadData, payloadError := payload.GetData(execution)

	if payloadError != nil {
		err := executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             collectPayloadDataStep.ID,
			ActionMessages: []string{"Failed to get Payload Data"},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			executions.EndWithError(execution)
			return
		}

		executions.EndWithError(execution)
		return
	}

	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             collectPayloadDataStep.ID,
		ActionMessages: []string{"Payload Data received"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	if flowDataErr == nil && payloadError == nil {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             collectDataStep.ID,
			ActionMessages: []string{"Data collected"},
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			executions.EndWithError(execution)
			return
		}
	} else {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             collectDataStep.ID,
			ActionMessages: []string{"Data collection finished with errors"},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			executions.EndWithError(execution)
			return
		}
	}

	// check for patterns
	match, err := flow.CheckPatterns(flowData, execution, payloadData)
	if err != nil {
		log.Error(err)
		executions.EndWithError(execution)
		return
	}

	if !match {
		return
	}

	// check for flow actions
	status, err := flow.CheckFlowActions(flowData, execution)
	if err != nil {
		log.Error(err)
		executions.EndWithError(execution)
		return
	}

	if !status {
		return
	}

	// set global variables
	variables.CurrentExecution = execution
	variables.CurrentPayload = payloadData
	variables.CurrentFlowActions = flowData.Actions

	// start every defined flow action
	// var actionsErrored []string
	if flowData.ExecParallel {
		for _, action := range flowData.Actions {
			if action.Active {
				go func(action models.Actions, execution models.Execution) {
					status, err := actions.StartAction(action, execution)
					if err != nil {
						log.Error(err)
						executions.EndWithError(execution)
						return
					}

					if !status {
						executions.EndWithError(execution)
						return
					}
				}(action, execution)
			}
		}
	} else {
		for _, action := range flowData.Actions {
			if action.Active {
				status, err := actions.StartAction(action, execution)
				if err != nil {
					log.Error(err)
					executions.EndWithError(execution)
					return
				}

				if !status {
					executions.EndWithError(execution)
					return
				}
			}
		}
	}

	executions.EndSuccess(execution)
}
