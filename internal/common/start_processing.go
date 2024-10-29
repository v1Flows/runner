package common

import (
	"alertflow-runner/config"
	"alertflow-runner/internal/executions"
	"alertflow-runner/internal/runner"
	"alertflow-runner/pkg/models"
	"time"

	log "github.com/sirupsen/logrus"
)

func startProcessing(execution models.Execution) {
	// ensure that runnerID is empty or equal to the current runnerID
	if execution.RunnerID != "" && execution.RunnerID != config.Config.RunnerID {
		log.Warnf("Execution %s is already picked up by another runner", execution.ID)
		return
	}

	execution.RunnerID = config.Config.RunnerID
	execution.Pending = false
	execution.Running = true
	execution.ExecutedAt = time.Now()

	err := executions.Update(execution)
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// set runner to busy
	runner.Busy(true)

	// send initial step to alertflow
	initialSteps, err := executions.SendInitialSteps(execution)
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// process each initial step where pending is true
	var flow models.Flows
	var payload models.Payload
	for _, step := range initialSteps {
		if step.Pending {
			data, _, canceled, no_pattern_match, failed, err := processStep(flow, payload, initialSteps, step, execution)
			if err != nil {
				executions.EndWithError(execution)
				return
			}

			if data["flow"] != nil {
				flow = data["flow"].(models.Flows)
			}

			if data["payload"] != nil {
				payload = data["payload"].(models.Payload)
			}

			if failed {
				executions.CancelRemainingSteps(execution.ID.String())
				executions.EndWithError(execution)
				return
			} else if canceled {
				executions.CancelRemainingSteps(execution.ID.String())
				executions.EndCanceled(execution)
				return
			} else if no_pattern_match {
				executions.CancelRemainingSteps(execution.ID.String())
				executions.EndNoPatternMatch(execution)
				return
			}
		}
	}

	// send flow actions as steps to alertflow
	flowActionStepsWithIDs, err := executions.SendFlowActionSteps(execution, flow)
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	// process each flow action step where pending is true
	for _, step := range flowActionStepsWithIDs {
		if step.Pending {
			_, _, canceled, no_pattern_match, failed, err := processStep(flow, payload, flowActionStepsWithIDs, step, execution)
			if err != nil {
				executions.EndWithError(execution)
				return
			}

			if failed {
				executions.CancelRemainingSteps(execution.ID.String())
				executions.EndWithError(execution)
				return
			} else if canceled {
				executions.CancelRemainingSteps(execution.ID.String())
				executions.EndCanceled(execution)
				return
			} else if no_pattern_match {
				executions.CancelRemainingSteps(execution.ID.String())
				executions.EndNoPatternMatch(execution)
				return
			}
		}
	}

	executions.EndSuccess(execution)

	runner.Busy(false)
}
