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
			data, _, canceled, failed, err := processStep(flow, payload, initialSteps, step, execution)
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
			}
		}
	}

	executions.EndSuccess(execution)

	runner.Busy(false)

	// var actionsToRun []string
	// var actionsRunStarted []string
	// var actionsRunFinished []string
	// var actionsRunCancelled []string
	// var actionsRunFailed []string

	// // start every defined flow action
	// if flowData.ExecParallel {
	// 	for _, action := range flowData.Actions {
	// 		if action.Active {
	// 			actionsToRun = append(actionsToRun, action.Name)

	// 			go func(action models.Actions, execution models.Execution) {
	// 				finished, canceled, failed, err := actions.StartAction(action, execution)
	// 				if err != nil {
	// 					log.Error(err)
	// 					executions.EndWithError(execution)
	// 					return
	// 				}

	// 				actionsRunStarted = append(actionsRunStarted, action.Name)

	// 				if failed {
	// 					actionsRunFailed = append(actionsRunFailed, action.Name)
	// 					return
	// 				} else if canceled {
	// 					actionsRunCancelled = append(actionsRunCancelled, action.Name)
	// 					return
	// 				} else if finished {
	// 					actionsRunFinished = append(actionsRunFinished, action.Name)
	// 				}
	// 			}(action, execution)
	// 		}
	// 	}

	// 	// wait for all actions to finish
	// 	for {
	// 		if len(actionsToRun) == len(actionsRunStarted) {
	// 			break
	// 		}

	// 		time.Sleep(1 * time.Second)
	// 	}
	// } else {
	// 	for _, action := range flowData.Actions {
	// 		if action.Active {
	// 			finished, canceled, failed, err := actions.StartAction(action, execution)
	// 			if err != nil {
	// 				log.Error(err)
	// 				executions.EndWithError(execution)
	// 				return
	// 			}

	// 			if failed {
	// 				actionsRunFailed = append(actionsRunFailed, action.Name)
	// 				executions.EndWithError(execution)
	// 				return
	// 			} else if canceled {
	// 				actionsRunCancelled = append(actionsRunCancelled, action.Name)
	// 				executions.EndCanceled(execution)
	// 				return
	// 			} else if finished {
	// 				actionsRunFinished = append(actionsRunFinished, action.Name)
	// 			}
	// 		}
	// 	}
	// }

	// if len(actionsRunFailed) > 0 {
	// 	executions.EndWithError(execution)
	// } else if len(actionsRunCancelled) > 0 {
	// 	executions.EndCanceled(execution)
	// } else {
	// 	executions.EndSuccess(execution)
	// }
}
