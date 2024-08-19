package processing

import (
	"alertflow-runner/functions/executions"
	"alertflow-runner/functions/flow"
	"alertflow-runner/functions/payload"
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"fmt"
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
	execution.TotalSteps = 2

	// update execution
	executions.Update(execution)

	// set runner picked up step
	stepData, _ := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID: execution.ID.String(),
		ActionName:  "Runner Pick Up",
		ActionMessages: []string{
			"Waiting for Runner to pick up Execution",
			"Runner picked up execution",
		},
		Finished:   true,
		FinishedAt: time.Now(),
	})

	// get flow data
	stepData, _ = executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Get Flow Data",
		ActionMessages: []string{"Requesting Flow Data from API"},
		StartedAt:      time.Now(),
	})

	flowData, flowActionData, err := flow.GetFlowData(execution)

	fmt.Println("Flow: ", flowData)

	if err != nil {
		executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             stepData.ID,
			ActionMessages: []string{"Failed to get Flow Data from API"},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})

		execution.FinishedAt = time.Now()
		execution.Running = false
		execution.Error = true
		executions.End(execution)
		return
	} else {
		executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             stepData.ID,
			ActionMessages: []string{"Flow Data received from API"},
			Finished:       true,
			FinishedAt:     time.Now(),
		})
	}

	// get payload data
	stepData, _ = executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Get Payload Data",
		ActionMessages: []string{"Requesting Payload Data from API"},
		StartedAt:      time.Now(),
	})

	payloadData, payloadError := payload.GetData(execution)

	fmt.Println("Payload: ", payloadData)

	if payloadError != nil {
		executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             stepData.ID,
			ActionMessages: []string{"Failed to get Payload Data from API"},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})

		execution.FinishedAt = time.Now()
		execution.Running = false
		execution.Error = true
		executions.End(execution)
		return
	} else {
		executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             stepData.ID,
			ActionMessages: []string{"Payload Data received from API"},
			Finished:       true,
			FinishedAt:     time.Now(),
		})
	}

	// check for flow actions
	stepData, _ = executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Check for Actions",
		ActionMessages: []string{"Checking if Flow has any Actions"},
		StartedAt:      time.Now(),
	})

	status := flow.CheckFlowActions(flowActionData)

	if !status {
		executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             stepData.ID,
			ActionMessages: []string{"Flow has no Actions defined"},
			NoResult:       true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})

		execution.FinishedAt = time.Now()
		execution.Running = false
		execution.Ghost = true
		executions.End(execution)
		return
	} else {
		executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             stepData.ID,
			ActionMessages: []string{"Flow has Actions defined"},
			Finished:       true,
			FinishedAt:     time.Now(),
		})
	}
}
