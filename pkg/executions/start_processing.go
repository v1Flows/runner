package executions

import (
	"time"

	"github.com/google/uuid"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/runner"
	"github.com/v1Flows/runner/pkg/plugins"

	log "github.com/sirupsen/logrus"
)

func startProcessing(platform string, cfg config.Config, actions []models.Actions, loadedPlugins map[string]plugins.Plugin, execution bmodels.Executions) {
	configManager := config.GetInstance()

	// ensure that execution runnerid equals the config runnerid
	if execution.RunnerID != configManager.GetRunnerID(platform) {
		log.Warnf("Execution %s is already picked up by another runner", execution.ID)
		return
	}

	execution.Status = "running"
	execution.ExecutedAt = time.Now()

	err := Update(cfg, execution)
	if err != nil {
		EndWithError(cfg, execution)
		return
	}

	// set runner to busy
	runner.Busy(platform, cfg, true)

	// send initial step to alertflow
	initialSteps, err := sendInitialSteps(cfg, actions, execution)
	if err != nil {
		EndWithError(cfg, execution)
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
				CancelRemainingSteps(cfg, execution.ID.String())
				// end execution
				EndWithError(cfg, execution)
				return
			}

			if res.Flow != nil {
				flow = *res.Flow
			} else if flow.ID == uuid.Nil {
				log.Error("Error parsing flow")
				CancelRemainingSteps(cfg, execution.ID.String())
				EndWithError(cfg, execution)
				return
			}

			if res.Alert != nil {
				alert = *res.Alert
			} else if alert.ID == uuid.Nil {
				log.Error("Error parsing alert")
				CancelRemainingSteps(cfg, execution.ID.String())
				EndWithError(cfg, execution)
				return
			}

			if res.Data["status"] == "noPatternMatch" {
				CancelRemainingSteps(cfg, execution.ID.String())
				EndNoPatternMatch(cfg, execution)
				return
			}

			if res.Data["status"] == "canceled" {
				CancelRemainingSteps(cfg, execution.ID.String())
				EndCanceled(cfg, execution)
				return
			}

			if !success {
				CancelRemainingSteps(cfg, execution.ID.String())
				EndWithError(cfg, execution)
				return
			}
		}
	}

	// send flow actions as steps to alertflow
	flowActionStepsWithIDs, err := sendFlowActionSteps(cfg, execution, flow)
	if err != nil {
		EndWithError(cfg, execution)
		return
	}

	if !flow.ExecParallel {
		// process each flow action step in sequential order where pending is true
		for _, step := range flowActionStepsWithIDs {
			if step.Status == "pending" {
				res, success, err := processStep(cfg, actions, loadedPlugins, flow, alert, flowActionStepsWithIDs, step, execution)
				if err != nil {
					// cancel remaining steps
					CancelRemainingSteps(cfg, execution.ID.String())
					// end execution
					EndWithError(cfg, execution)
					return
				}

				if res.Data["status"] == "noPatternMatch" {
					CancelRemainingSteps(cfg, execution.ID.String())
					EndNoPatternMatch(cfg, execution)
					return
				}

				if res.Data["status"] == "canceled" {
					CancelRemainingSteps(cfg, execution.ID.String())
					EndCanceled(cfg, execution)
					return
				}

				if !success {
					CancelRemainingSteps(cfg, execution.ID.String())
					EndWithError(cfg, execution)
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
			EndWithError(cfg, execution)
			return
		}

		if canceledSteps > 0 {
			EndCanceled(cfg, execution)
			return
		}

		if noPatternMatchSteps > 0 {
			EndNoPatternMatch(cfg, execution)
			return
		}
	}

	EndSuccess(cfg, execution)

	runner.Busy(platform, cfg, false)
}
