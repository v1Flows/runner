package processing

import (
	"alertflow-runner/functions/actions"
	"alertflow-runner/functions/executions"
	"alertflow-runner/functions/flow"
	"alertflow-runner/functions/payload"
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"time"
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
	stepData, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Runner Pick Up",
		ActionMessages: []string{"Waiting for Runner to pick up Execution", "Runner picked up execution"},
		Finished:       true,
		FinishedAt:     time.Now(),
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
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	_, flowActionData, flowDataErr := flow.GetFlowData(execution)

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
			ActionMessages: []string{"Collecting Data finished"},
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
			ActionMessages: []string{"Collecting Data finished with errors"},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			executions.EndWithError(execution)
			return
		}
	}

	// check for flow actions
	stepData, err = executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Check for Actions",
		ActionMessages: []string{"Checking if Flow has any Actions"},
		StartedAt:      time.Now(),
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	status := flow.CheckFlowActions(flowActionData)

	if !status {
		err := executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             stepData.ID,
			ActionMessages: []string{"Flow has no Actions defined"},
			NoResult:       true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			executions.EndWithError(execution)
			return
		}

		executions.EndWithGhost(execution)
		return
	}

	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             stepData.ID,
		ActionMessages: []string{"Actions found in Flow"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// start every defined flow action
	var actionsFinished []string
	var actionsNoMatch []string
	for _, action := range flowActionData {
		if action.Status {
			success, no_match, err := actions.StartAction(execution, action, payloadData)
			if err != nil {
				break
			}

			if success {
				actionsFinished = append(actionsFinished, action.Name)
			} else if no_match {
				actionsNoMatch = append(actionsNoMatch, action.Name)
			}
		}
	}

	if len(actionsFinished) > 0 {
		executions.EndSuccess(execution)
		return
	} else if len(actionsFinished) == 0 && len(actionsNoMatch) > 0 {
		executions.EndWithNoMatch(execution)
		return
	} else {
		executions.EndWithError(execution)
		return
	}
}
