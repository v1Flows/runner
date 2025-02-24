package common

import (
	"time"

	"github.com/AlertFlow/runner/config"
	internal_executions "github.com/AlertFlow/runner/internal/executions"
	"github.com/AlertFlow/runner/internal/runner"
	"github.com/AlertFlow/runner/pkg/executions"
	"github.com/AlertFlow/runner/pkg/plugins"
	"github.com/google/uuid"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

func startProcessing(cfg config.Config, actions []models.Actions, loadedPlugins map[string]plugins.Plugin, execution bmodels.Executions) {
	// ensure that execution runnerid equals the config runnerid
	if execution.RunnerID != cfg.Alertflow.RunnerID {
		log.Warnf("Execution %s is already picked up by another runner", execution.ID)
		return
	}

	execution.Status = "running"
	execution.ExecutedAt = time.Now()

	err := executions.Update(cfg, execution)
	if err != nil {
		executions.EndWithError(cfg, execution)
		return
	}

	// set runner to busy
	runner.Busy(cfg, true)

	// send initial step to alertflow
	initialSteps, err := internal_executions.SendInitialSteps(cfg, actions, execution)
	if err != nil {
		executions.EndWithError(cfg, execution)
		return
	}

	// process each initial step where pending is true
	var flow bmodels.Flows
	var alert bmodels.Alerts
	for _, step := range initialSteps {
		if step.Status == "pending" {
			res, success, err := processStep(cfg, actions, loadedPlugins, flow, alert, initialSteps, step, execution)
			if err != nil {
				log.Error("Error processing initial step: ", err)
				// cancel remaining steps
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				// end execution
				executions.EndWithError(cfg, execution)
				return
			}

			if res.Flow != nil {
				flow = *res.Flow
			} else if flow.ID == uuid.Nil {
				log.Error("Error parsing flow")
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution)
				return
			}

			if res.Alert != nil {
				alert = *res.Alert
			} else if alert.ID == uuid.Nil {
				log.Error("Error parsing alert")
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution)
				return
			}

			if res.Data["status"] == "noPatternMatch" {
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndNoPatternMatch(cfg, execution)
				return
			}

			if res.Data["status"] == "canceled" {
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndCanceled(cfg, execution)
				return
			}

			if !success {
				internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution)
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
				res, success, err := processStep(cfg, actions, loadedPlugins, flow, alert, flowActionStepsWithIDs, step, execution)
				if err != nil {
					// cancel remaining steps
					internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
					// end execution
					executions.EndWithError(cfg, execution)
					return
				}

				if res.Data["status"] == "noPatternMatch" {
					internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
					executions.EndNoPatternMatch(cfg, execution)
					return
				}

				if res.Data["status"] == "canceled" {
					internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
					executions.EndCanceled(cfg, execution)
					return
				}

				if !success {
					internal_executions.CancelRemainingSteps(cfg, execution.ID.String())
					executions.EndWithError(cfg, execution)
					return
				}
			}
		}
	} else {
		var executedSteps int
		var failedSteps int
		var noPatternMatchSteps int
		var canceledSteps int
		var successSteps int
		// process each flow action step in parallel where pending is true
		for _, step := range flowActionStepsWithIDs {
			if step.Status == "pending" {
				go func() {
					res, success, err := processStep(cfg, actions, loadedPlugins, flow, alert, flowActionStepsWithIDs, step, execution)
					if err != nil {
						failedSteps++
					}

					executedSteps++

					if res.Data["status"] == "noPatternMatch" {
						noPatternMatchSteps++
					}

					if res.Data["status"] == "canceled" {
						canceledSteps++
					}

					if !success {
						failedSteps++
					}

					if success {
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
		}

		if canceledSteps > 0 {
			executions.EndCanceled(cfg, execution)
			return
		}

		if noPatternMatchSteps > 0 {
			executions.EndNoPatternMatch(cfg, execution)
			return
		}
	}

	executions.EndSuccess(cfg, execution)

	runner.Busy(cfg, false)
}
