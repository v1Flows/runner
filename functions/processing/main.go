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
	stepData, _ := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:   execution.ID.String(),
		ActionName:    "Get Flow Data",
		ActionMessage: "Requesting Flow Data from API",
		StartedAt:     time.Now(),
	})

	flowData, _ := flow.GetFlowData(execution)

	time.Sleep(15 * time.Second)

	executions.UpdateStep(execution, models.ExecutionSteps{
		ID:         stepData.ID,
		Finished:   true,
		FinishedAt: time.Now(),
	})

	// check for flow actions
	stepData, _ = executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:   execution.ID.String(),
		ActionName:    "Check for Actions",
		ActionMessage: "Checking if Flow has any Actions",
		StartedAt:     time.Now(),
	})

	status := flow.CheckFlowActions(flowData)

	time.Sleep(15 * time.Second)

	if !status {
		executions.UpdateStep(execution, models.ExecutionSteps{
			ID:         stepData.ID,
			Finished:   true,
			FinishedAt: time.Now(),
		})

		execution.FinishedAt = time.Now()
		execution.Running = false
		executions.End(execution)
		return
	}
}
