package exflow

import (
	"time"

	"github.com/google/uuid"
	af_models "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/runner"
	"github.com/v1Flows/runner/pkg/executions"
	"github.com/v1Flows/runner/pkg/plugins"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func StartProcessing(platform string, cfg config.Config, actions []shared_models.Actions, loadedPlugins map[string]plugins.Plugin, execution shared_models.Executions) {
	configManager := config.GetInstance()

	// ensure that execution runnerid equals the config runnerid
	if execution.RunnerID != configManager.GetRunnerID(platform) {
		log.Warnf("Execution %s is already picked up by another runner", execution.ID)
		return
	}

	execution.Status = "running"
	execution.ExecutedAt = time.Now()

	err := executions.UpdateExecution(cfg, execution)
	if err != nil {
		executions.EndWithError(cfg, execution)
		return
	}

	// set runner to busy
	runner.Busy(platform, cfg, true)

	// send initial step to alertflow
	initialSteps, err := sendInitialSteps(cfg, actions, execution)
	if err != nil {
		executions.EndWithError(cfg, execution)
		return
	}

	// process each initial step where pending is true
	var flow shared_models.Flows
	for _, step := range initialSteps {
		if step.Status == "pending" {
			res, success, err := executions.ProcessStep(cfg, actions, loadedPlugins, flow, af_models.Alerts{}, initialSteps, step, execution)
			if err != nil {
				log.Error("Error processing initial step: ", err)
				// cancel remaining steps
				executions.CancelRemainingSteps(cfg, execution.ID.String())
				// end execution
				executions.EndWithError(cfg, execution)
				return
			}

			if res.Flow != nil {
				flow = *res.Flow
			} else if flow.ID == uuid.Nil {
				log.Error("Error parsing flow")
				executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution)
				return
			}

			if res.Data["status"] == "noPatternMatch" {
				executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndNoPatternMatch(cfg, execution)
				return
			}

			if res.Data["status"] == "canceled" {
				executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndCanceled(cfg, execution)
				return
			}

			if !success {
				executions.CancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution)
				return
			}
		}
	}

	// send flow actions as steps to exflow
	flowActionStepsWithIDs, err := executions.SendFlowActionSteps(cfg, execution, flow)
	if err != nil {
		executions.EndWithError(cfg, execution)
		return
	}

	if !flow.ExecParallel {
		// process each flow action step in sequential order where pending is true
		for _, step := range flowActionStepsWithIDs {
			if step.Status == "pending" {
				res, success, err := executions.ProcessStep(cfg, actions, loadedPlugins, flow, af_models.Alerts{}, flowActionStepsWithIDs, step, execution)
				if err != nil {
					// cancel remaining steps
					executions.CancelRemainingSteps(cfg, execution.ID.String())
					// end execution
					executions.EndWithError(cfg, execution)
					return
				}

				if res.Data["status"] == "noPatternMatch" {
					executions.CancelRemainingSteps(cfg, execution.ID.String())
					executions.EndNoPatternMatch(cfg, execution)
					return
				}

				if res.Data["status"] == "canceled" {
					executions.CancelRemainingSteps(cfg, execution.ID.String())
					executions.EndCanceled(cfg, execution)
					return
				}

				if !success {
					executions.CancelRemainingSteps(cfg, execution.ID.String())
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
					res, success, err := executions.ProcessStep(cfg, actions, loadedPlugins, flow, af_models.Alerts{}, flowActionStepsWithIDs, step, execution)
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

	runner.Busy(platform, cfg, false)
}
