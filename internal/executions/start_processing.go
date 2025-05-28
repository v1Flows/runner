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

func startProcessing(platform string, actions []shared_models.Action, loadedPlugins map[string]plugins.Plugin, execution shared_models.Executions, alertID string) {
	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	// ensure that execution runnerid equals the config runnerid
	if execution.RunnerID != configManager.GetRunnerID(platform) {
		log.Warnf("Execution %s is already picked up by another runner", execution.ID)
		return
	}

	// Start sending heartbeats
	doneHeartbeat := make(chan struct{})
	go func() {
		for {
			select {
			case <-doneHeartbeat:
				// Stop sending heartbeats when done
				return
			default:
				// Send a heartbeat
				executions.SendHeartbeat(nil, execution, platform)
				// Wait for a short interval before sending the next heartbeat
				time.Sleep(5 * time.Second) // Adjust the interval as needed
			}
		}
	}()

	// create workspace dir for execution
	workspace := fmt.Sprintf("%s/%s", cfg.WorkspaceDir, execution.ID)
	err := os.MkdirAll(workspace, 0755)
	if err != nil {
		log.Error("Error creating workspace dir: ", err)
	}

	execution.Status = "running"
	execution.ExecutedAt = time.Now()

	err = executions.UpdateExecution(nil, execution, platform)
	if err != nil {
		executions.EndWithError(nil, execution, platform)
		// Stop heartbeats and finish processing
		close(doneHeartbeat)
		finishProcessing(platform, cfg, execution)
		return
	}

	// set runner to busy
	runner.Busy(platform, true)

	// send initial step
	var initialSteps []shared_models.ExecutionSteps
	if platform == "alertflow" {
		initialSteps, err = internal_alertflow.SendInitialSteps(cfg, actions, execution, alertID)
		if err != nil {
			executions.EndWithError(nil, execution, platform)
			// Stop heartbeats and finish processing
			close(doneHeartbeat)
			finishProcessing(platform, cfg, execution)
			return
		}
	} else if platform == "exflow" {
		initialSteps, err = internal_exflow.SendInitialSteps(cfg, actions, execution)
		if err != nil {
			executions.EndWithError(nil, execution, platform)
			// Stop heartbeats and finish processing
			close(doneHeartbeat)
			finishProcessing(platform, cfg, execution)
			return
		}
	}

	// process each initial step where pending is true
	var flow shared_models.Flows
	var flowBytes []byte
	var alert bmodels.Alerts
	for _, step := range initialSteps {
		if step.Status == "pending" {
			res, success, canceled, err := processStep(cfg, workspace, actions, loadedPlugins, flow, flowBytes, alert, initialSteps, step, execution)
			if err != nil {
				log.Error("Error processing initial step: ", err)
				// cancel remaining steps
				cancelRemainingSteps(execution.ID.String())
				// end execution
				executions.EndWithError(nil, execution, platform)
				// Stop heartbeats and finish processing
				close(doneHeartbeat)
				finishProcessing(platform, cfg, execution)
				return
			}

			if canceled {
				cancelRemainingSteps(execution.ID.String())
				executions.EndCanceled(nil, execution, platform)
				// Stop heartbeats and finish processing
				close(doneHeartbeat)
				finishProcessing(platform, cfg, execution)
				return
			}

			if res.Flow != nil {
				flow = *res.Flow
			} else if flow.ID == uuid.Nil {
				log.Error("Error parsing flow")
				cancelRemainingSteps(execution.ID.String())
				executions.EndWithError(nil, execution, platform)
				// Stop heartbeats and finish processing
				close(doneHeartbeat)
				finishProcessing(platform, cfg, execution)
				return
			}

			if len(res.FlowBytes) > 0 {
				flowBytes = res.FlowBytes
			}

			if platform == "alertflow" && res.Alert != nil {
				alert = *res.Alert
			} else if platform == "alertflow" && alert.ID == uuid.Nil {
				log.Error("Error parsing alert")
				cancelRemainingSteps(execution.ID.String())
				executions.EndWithError(nil, execution, platform)
				// Stop heartbeats and finish processing
				close(doneHeartbeat)
				finishProcessing(platform, cfg, execution)
				return
			}

			if res.Data["status"] == "noPatternMatch" {
				cancelRemainingSteps(execution.ID.String())
				executions.EndNoPatternMatch(nil, execution, platform)
				finishProcessing(platform, cfg, execution)
				return
			}

			if res.Data["status"] == "canceled" {
				cancelRemainingSteps(execution.ID.String())
				executions.EndCanceled(nil, execution, platform)
				// Stop heartbeats and finish processing
				close(doneHeartbeat)
				finishProcessing(platform, cfg, execution)
				return
			}

			if !success {
				cancelRemainingSteps(execution.ID.String())
				executions.EndWithError(nil, execution, platform)
				// Stop heartbeats and finish processing
				close(doneHeartbeat)
				finishProcessing(platform, cfg, execution)
				return
			}
		}
	}

	// send flow actions as steps to alertflow
	flowActionStepsWithIDs, err := sendFlowActionSteps(cfg, execution, flow)
	if err != nil {
		executions.EndWithError(nil, execution, platform)
		// Stop heartbeats and finish processing
		close(doneHeartbeat)
		finishProcessing(platform, cfg, execution)
		return
	}

	if !flow.ExecParallel {
		// process each flow action step in sequential order where pending is true
		for _, step := range flowActionStepsWithIDs {
			if step.Status == "pending" {
				res, success, canceled, err := processStep(cfg, workspace, actions, loadedPlugins, flow, flowBytes, alert, flowActionStepsWithIDs, step, execution)
				if err != nil {
					// cancel remaining steps
					cancelRemainingSteps(execution.ID.String())

					// start failure pipeline
					if flow.FailurePipelineID != "" || step.Action.FailurePipelineID != "" {
						err = startFailurePipeline(cfg, workspace, actions, loadedPlugins, flow, flowBytes, alert, flowActionStepsWithIDs, step, execution)
						if err != nil {
							// end execution with recovered status
							executions.EndWithError(nil, execution, platform)
							// Stop heartbeats and finish processing
							close(doneHeartbeat)
							finishProcessing(platform, cfg, execution)
							return
						}
					}

					// end execution
					executions.EndWithRecovered(nil, execution, platform)
					// Stop heartbeats and finish processing
					close(doneHeartbeat)
					finishProcessing(platform, cfg, execution)
					return
				}

				if res.Data["status"] == "noPatternMatch" {
					cancelRemainingSteps(execution.ID.String())
					executions.EndNoPatternMatch(nil, execution, platform)
					// Stop heartbeats and finish processing
					close(doneHeartbeat)
					finishProcessing(platform, cfg, execution)
					return
				}

				if res.Data["status"] == "canceled" {
					cancelRemainingSteps(execution.ID.String())
					executions.EndCanceled(nil, execution, platform)
					// Stop heartbeats and finish processing
					close(doneHeartbeat)
					finishProcessing(platform, cfg, execution)
					return
				}

				if canceled {
					cancelRemainingSteps(execution.ID.String())
					executions.EndCanceled(nil, execution, platform)
					// Stop heartbeats and finish processing
					close(doneHeartbeat)
					finishProcessing(platform, cfg, execution)
					return
				}

				if !success {
					cancelRemainingSteps(execution.ID.String())

					// start failure pipeline if enabled
					if flow.FailurePipelineID != "" || step.Action.FailurePipelineID != "" {
						err = startFailurePipeline(cfg, workspace, actions, loadedPlugins, flow, flowBytes, alert, flowActionStepsWithIDs, step, execution)
						if err != nil {
							executions.EndWithError(nil, execution, platform)
							// Stop heartbeats and finish processing
							close(doneHeartbeat)
							finishProcessing(platform, cfg, execution)
							return
						}

						// end execution with recovered status
						executions.EndWithRecovered(nil, execution, platform)
						// Stop heartbeats and finish processing
						close(doneHeartbeat)
						finishProcessing(platform, cfg, execution)
						return
					}

					executions.EndWithError(nil, execution, platform)
					// Stop heartbeats and finish processing
					close(doneHeartbeat)
					finishProcessing(platform, cfg, execution)
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
					res, success, canceled, err := processStep(cfg, workspace, actions, loadedPlugins, flow, flowBytes, alert, flowActionStepsWithIDs, step, execution)
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

					if canceled {
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
			executions.EndWithError(nil, execution, platform)
			// Stop heartbeats and finish processing
			close(doneHeartbeat)
			finishProcessing(platform, cfg, execution)
			return
		}

		if canceledSteps > 0 {
			executions.EndCanceled(nil, execution, platform)
			// Stop heartbeats and finish processing
			close(doneHeartbeat)
			finishProcessing(platform, cfg, execution)
			return
		}

		if noPatternMatchSteps > 0 {
			executions.EndNoPatternMatch(nil, execution, platform)
			// Stop heartbeats and finish processing
			close(doneHeartbeat)
			finishProcessing(platform, cfg, execution)
			return
		}
	}

	executions.EndSuccess(nil, execution, platform)

	// Stop heartbeats and finish processing
	close(doneHeartbeat)
	finishProcessing(platform, cfg, execution)
}

func finishProcessing(platform string, cfg config.Config, execution shared_models.Executions) {
	err := os.RemoveAll(fmt.Sprintf("%s/%s", cfg.WorkspaceDir, execution.ID))
	if err != nil {
		log.Error("Error deleting workspace dir: ", err)
	}

	runner.Busy(platform, false)
}
