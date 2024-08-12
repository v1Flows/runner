package processing

import (
	"alertflow-runner/functions/executions"
	"alertflow-runner/functions/flow"
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
	executions.Update(execution)

	// get flow data
	executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:   execution.ID.String(),
		ActionName:    "Get Flow Data",
		ActionMessage: "Requesting Flow Data from API",
		StartedAt:     time.Now(),
	})

	flowData, _ := flow.GetFlowData(execution)

	executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:   execution.ID.String(),
		ActionName:    "Get Flow Data",
		ActionMessage: "Requesting Flow Data from API finished",
		Finished:      true,
		FinishedAt:    time.Now(),
	})

	// check for flow actions
	executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:   execution.ID.String(),
		ActionName:    "Check for Actions",
		ActionMessage: "Checking if Flow has any Actions",
		StartedAt:     time.Now(),
	})

	status := flow.CheckFlowActions(flowData)

	if !status {
		executions.SendStep(execution, models.ExecutionSteps{
			ExecutionID:   execution.ID.String(),
			ActionName:    "Check for Actions",
			ActionMessage: "No Flow Actions found",
			Finished:      true,
			FinishedAt:    time.Now(),
		})

		execution.FinishedAt = time.Now()
		execution.Running = false
		executions.End(execution)
		return
	}
}
