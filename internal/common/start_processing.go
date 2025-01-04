package common

import (
	"time"

	"gitlab.justlab.xyz/alertflow-public/runner/config"
	internal_executions "gitlab.justlab.xyz/alertflow-public/runner/internal/executions"
	"gitlab.justlab.xyz/alertflow-public/runner/internal/runner"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/executions"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func startProcessing(execution models.Execution) {
	// ensure that runnerID is empty or equal to the current runnerID
	if execution.RunnerID != "" && execution.RunnerID != config.Config.Alertflow.RunnerID {
		log.Warnf("Execution %s is already picked up by another runner", execution.ID)
		return
	}

	execution.RunnerID = config.Config.Alertflow.RunnerID
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
	initialSteps, err := internal_executions.SendInitialSteps(execution)
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
				internal_executions.CancelRemainingSteps(execution.ID.String())
				executions.EndWithError(execution)
				return
			} else if canceled {
				internal_executions.CancelRemainingSteps(execution.ID.String())
				executions.EndCanceled(execution)
				return
			} else if no_pattern_match {
				internal_executions.CancelRemainingSteps(execution.ID.String())
				executions.EndNoPatternMatch(execution)
				return
			}
		}
	}

	// send flow actions as steps to alertflow
	flowActionStepsWithIDs, err := internal_executions.SendFlowActionSteps(execution, flow)
	if err != nil {
		executions.EndWithError(execution)
		return
	}

	if !flow.ExecParallel {
		// process each flow action step in sequential order where pending is true
		for _, step := range flowActionStepsWithIDs {
			if step.Pending {
				_, _, canceled, no_pattern_match, failed, err := processStep(flow, payload, flowActionStepsWithIDs, step, execution)
				if err != nil {
					executions.EndWithError(execution)
					return
				}

				if failed {
					internal_executions.CancelRemainingSteps(execution.ID.String())
					executions.EndWithError(execution)
					return
				} else if canceled {
					internal_executions.CancelRemainingSteps(execution.ID.String())
					executions.EndCanceled(execution)
					return
				} else if no_pattern_match {
					internal_executions.CancelRemainingSteps(execution.ID.String())
					executions.EndNoPatternMatch(execution)
					return
				}
			}
		}
	} else {
		var executedSteps int
		var failedSteps int
		var canceledSteps int
		var noPatternMatchSteps int
		var successSteps int
		// process each flow action step in parallel where pending is true
		for _, step := range flowActionStepsWithIDs {
			if step.Pending {
				go func() {
					_, finished, cancleded, no_pattern_match, failed, err := processStep(flow, payload, flowActionStepsWithIDs, step, execution)
					if err != nil {
						executions.EndWithError(execution)
						return
					}

					executedSteps++

					if failed {
						failedSteps++
					} else if cancleded {
						canceledSteps++
					} else if no_pattern_match {
						noPatternMatchSteps++
					} else if finished {
						successSteps++
					}
				}()
			}
		}

		// wait for all steps to finish
		for executedSteps < len(flowActionStepsWithIDs) {
			if executedSteps == len(flowActionStepsWithIDs) {
				break
			}
		}

		if failedSteps > 0 {
			executions.EndWithError(execution)
			return
		} else if canceledSteps > 0 {
			executions.EndCanceled(execution)
			return
		} else if noPatternMatchSteps > 0 {
			executions.EndNoPatternMatch(execution)
			return
		}
	}

	executions.EndSuccess(execution)

	runner.Busy(false)
}
