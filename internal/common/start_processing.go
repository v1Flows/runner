package common

import (
	"time"

	"github.com/AlertFlow/runner/config"
	internal_executions "github.com/AlertFlow/runner/internal/executions"
	"github.com/AlertFlow/runner/internal/runner"
	"github.com/AlertFlow/runner/pkg/executions"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

func startProcessing(cfg config.Config, execution bmodels.Executions) {
	// ensure that runnerID is empty or equal to the current runnerID
	if execution.RunnerID != "" && execution.RunnerID != cfg.Alertflow.RunnerID {
		log.Warnf("Execution %s is already picked up by another runner", execution.ID)
		return
	}

	execution.RunnerID = cfg.Alertflow.RunnerID
	execution.Status = "running"
	execution.ExecutedAt = time.Now()

	err := executions.Update(cfg, execution)
	if err != nil {
		executions.EndWithError(cfg, execution)
		return
	}

	// set runner to busy
	runner.Busy(true)

	// send initial step to alertflow
	initialSteps, err := internal_executions.SendInitialSteps(cfg, execution)
	if err != nil {
		executions.EndWithError(cfg, execution)
		return
	}

	// process each initial step where pending is true
	var flow bmodels.Flows
	var payload bmodels.Payloads
	for _, step := range initialSteps {
		if step.Status == "pending" {
			data, _, canceled, no_pattern_match, failed, err := processStep(cfg, flow, payload, initialSteps, step, execution)
			if err != nil {
				executions.EndWithError(cfg, execution)
				return
			}

			if data["flow"] != nil {
				flow = data["flow"].(bmodels.Flows)
			}

			if data["payload"] != nil {
				payload = data["payload"].(bmodels.Payloads)
			}

			if failed {
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution)
				return
			} else if canceled {
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndCanceled(cfg, execution)
				return
			} else if no_pattern_match {
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndNoPatternMatch(cfg, execution)
				return
			}
		}
	}

	// send flow actions as steps to alertflow
	flowActionStepsWithIDs, err := internal_executions.SendFlowActionSteps(cfg, execution, flow)
	if err != nil {
		executions.EndWithError(cfg, execution)
		return
	}

	if !flow.ExecParallel {
		// process each flow action step in sequential order where pending is true
		for _, step := range flowActionStepsWithIDs {
			if step.Status == "pending" {
				_, _, canceled, no_pattern_match, failed, err := processStep(cfg, flow, payload, flowActionStepsWithIDs, step, execution)
				if err != nil {
					executions.EndWithError(cfg, execution)
					return
				}

				if failed {
					internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
					executions.EndWithError(cfg, execution)
					return
				} else if canceled {
					internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
					executions.EndCanceled(cfg, execution)
					return
				} else if no_pattern_match {
					internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
					executions.EndNoPatternMatch(cfg, execution)
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
			if step.Status == "pending" {
				go func() {
					_, finished, cancleded, no_pattern_match, failed, err := processStep(cfg, flow, payload, flowActionStepsWithIDs, step, execution)
					if err != nil {
						executions.EndWithError(cfg, execution)
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
			executions.EndWithError(cfg, execution)
			return
		} else if canceledSteps > 0 {
			executions.EndCanceled(cfg, execution)
			return
		} else if noPatternMatchSteps > 0 {
			executions.EndNoPatternMatch(cfg, execution)
			return
		}
	}

	executions.EndSuccess(cfg, execution)

	runner.Busy(false)
}
