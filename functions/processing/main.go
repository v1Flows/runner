package processing

import (
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"time"
)

func StartProcessing(execution models.Execution) {
	// set own runner id
	execution.RunnerID = config.Config.RunnerID
	// unset waiting
	execution.Waiting = false
	// set running
	execution.Running = true
	// set executed at
	execution.ExecutedAt = time.Now()

	// update execution
	UpdateExecution(execution)

	// get flow data
	SendExecutionStep(execution, models.ExecutionSteps{
		ExecutionID:   execution.ID.String(),
		ActionName:    "Get Flow Data",
		ActionMessage: "Requesting Flow Data from API",
		StartedAt:     time.Now(),
	})

	flow, _ := GetFlowData(execution)

	SendExecutionStep(execution, models.ExecutionSteps{
		ExecutionID:   execution.ID.String(),
		ActionName:    "Get Flow Data",
		ActionMessage: "Requesting Flow Data from API finished",
		Finished:      true,
		FinishedAt:    time.Now(),
	})

	// check for flow actions
	SendExecutionStep(execution, models.ExecutionSteps{
		ExecutionID:   execution.ID.String(),
		ActionName:    "Check for Actions",
		ActionMessage: "Checking if Flow has any Actions",
		StartedAt:     time.Now(),
	})

	status := CheckFlowActions(flow)

	if !status {
		SendExecutionStep(execution, models.ExecutionSteps{
			ExecutionID:   execution.ID.String(),
			ActionName:    "Check for Actions",
			ActionMessage: "No Flow Actions found",
			Finished:      true,
			FinishedAt:    time.Now(),
		})

		execution.FinishedAt = time.Now()
		execution.Running = false
		EndExecution(execution)
		return
	}
}
