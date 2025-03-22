package internal_executions

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	internal_alertflow "github.com/v1Flows/runner/internal/alertflow"
	internal_exflow "github.com/v1Flows/runner/internal/exflow"
	"github.com/v1Flows/runner/internal/runner"
	"github.com/v1Flows/runner/pkg/executions"
	"github.com/v1Flows/runner/pkg/plugins"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func startProcessing(platform string, cfg config.Config, actions []shared_models.Action, loadedPlugins map[string]plugins.Plugin, execution shared_models.Executions, alertID string) {
	configManager := config.GetInstance()

	// ensure that execution runnerid equals the config runnerid
	if execution.RunnerID != configManager.GetRunnerID(platform) {
		log.Warnf("Execution %s is already picked up by another runner", execution.ID)
		return
	}

	// create workspace dir for execution
	workspace := fmt.Sprintf("%s/%s", cfg.WorkspaceDir, execution.ID)
	err := os.MkdirAll(workspace, 0755)
	if err != nil {
		log.Error("Error creating workspace dir: ", err)
	}

	execution.Status = "running"
	execution.ExecutedAt = time.Now()

	err = executions.UpdateExecution(cfg, execution, platform)
	if err != nil {
		executions.EndWithError(cfg, execution, platform)
		return
	}

	// set runner to busy
	runner.Busy(platform, cfg, true)

	// send initial step
	var initialSteps []shared_models.ExecutionSteps
	if platform == "alertflow" {
		initialSteps, err = internal_alertflow.SendInitialSteps(cfg, actions, execution, alertID)
		if err != nil {
			executions.EndWithError(cfg, execution, platform)
			return
		}
	} else if platform == "exflow" {
		initialSteps, err = internal_exflow.SendInitialSteps(cfg, actions, execution)
		if err != nil {
			executions.EndWithError(cfg, execution, platform)
			return
		}
	}

	// process each initial step where pending is true
	var flow shared_models.Flows
	var flowBytes []byte
	var alert bmodels.Alerts
	for _, step := range initialSteps {
		if step.Status == "pending" {
			res, success, err := processStep(cfg, workspace, actions, loadedPlugins, flow, flowBytes, alert, initialSteps, step, execution)
			if err != nil {
				log.Error("Error processing initial step: ", err)
				// cancel remaining steps
				cancelRemainingSteps(cfg, execution.ID.String())
				// end execution
				executions.EndWithError(cfg, execution, platform)
				return
			}

			if res.Flow != nil {
				flow = *res.Flow
			} else if flow.ID == uuid.Nil {
				log.Error("Error parsing flow")
				cancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution, platform)
				return
			}

			if len(res.FlowBytes) > 0 {
				flowBytes = res.FlowBytes
			}

			if platform == "alertflow" && res.Alert != nil {
				alert = *res.Alert
			} else if platform == "alertflow" && alert.ID == uuid.Nil {
				log.Error("Error parsing alert")
				cancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution, platform)
				return
			}

			if res.Data["status"] == "noPatternMatch" {
				cancelRemainingSteps(cfg, execution.ID.String())
				executions.EndNoPatternMatch(cfg, execution, platform)
				return
			}

			if res.Data["status"] == "canceled" {
				cancelRemainingSteps(cfg, execution.ID.String())
				executions.EndCanceled(cfg, execution, platform)
				return
			}

			if !success {
				cancelRemainingSteps(cfg, execution.ID.String())
				executions.EndWithError(cfg, execution, platform)
				return
			}
		}
	}

	// send flow actions as steps to alertflow
	flowActionStepsWithIDs, err := sendFlowActionSteps(cfg, execution, flow)
	if err != nil {
		executions.EndWithError(cfg, execution, platform)
		return
	}

	if !flow.ExecParallel {
		// process each flow action step in sequential order where pending is true
		for _, step := range flowActionStepsWithIDs {
			if step.Status == "pending" {
				res, success, err := processStep(cfg, workspace, actions, loadedPlugins, flow, flowBytes, alert, flowActionStepsWithIDs, step, execution)
				if err != nil {
					// cancel remaining steps
					cancelRemainingSteps(cfg, execution.ID.String())
					// end execution
					executions.EndWithError(cfg, execution, platform)
					return
				}

				if res.Data["status"] == "noPatternMatch" {
					cancelRemainingSteps(cfg, execution.ID.String())
					executions.EndNoPatternMatch(cfg, execution, platform)
					return
				}

				if res.Data["status"] == "canceled" {
					cancelRemainingSteps(cfg, execution.ID.String())
					executions.EndCanceled(cfg, execution, platform)
					return
				}

				if !success {
					cancelRemainingSteps(cfg, execution.ID.String())
					executions.EndWithError(cfg, execution, platform)
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
					res, success, err := processStep(cfg, workspace, actions, loadedPlugins, flow, flowBytes, alert, flowActionStepsWithIDs, step, execution)
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
			executions.EndWithError(cfg, execution, platform)
			return
		}

		if canceledSteps > 0 {
			executions.EndCanceled(cfg, execution, platform)
			return
		}

		if noPatternMatchSteps > 0 {
			executions.EndNoPatternMatch(cfg, execution, platform)
			return
		}
	}

	executions.EndSuccess(cfg, execution, platform)

	err = os.RemoveAll(fmt.Sprintf("%s/%s", cfg.WorkspaceDir, execution.ID))
	if err != nil {
		log.Error("Error deleting workspace dir: ", err)
	}

	runner.Busy(platform, cfg, false)
}
