package common

import (
	"time"

	"github.com/AlertFlow/runner/internal/plugin"
	"github.com/AlertFlow/runner/pkg/executions"
	"github.com/AlertFlow/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

var actions []models.ActionDetails

func RegisterActions(loadedActions []models.ActionDetails) {
	actions = loadedActions
}

func processStep(pluginManager *plugin.Manager, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, execution models.Execution) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool, err error) {
	// set step to running
	step.Pending = false
	step.Running = true
	step.StartedAt = time.Now()
	step.RunnerID = execution.RunnerID

	if err := executions.UpdateStep(execution.ID.String(), step); err != nil {
		log.Error(err)
		return nil, false, false, false, false, err
	}

	valid, pluginVersion := checkActionVersionAgainstPluginVersion(step)

	if !valid {
		// dont execute step and quit execution
		step.ActionMessages = append(step.ActionMessages, "Action not compatible with plugin version", "Plugin Version: "+pluginVersion+" Action Version: "+step.ActionVersion, "Stopping execution")
		step.Running = false
		step.Error = true
		step.Finished = true
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(execution.ID.String(), step); err != nil {
			log.Error(err)
			return nil, false, false, false, false, err
		}

		return nil, false, false, false, true, nil
	}

	var found bool
	var action models.ActionDetails
	for _, a := range actions {
		if a.ID == step.ActionType {
			found = true
			action = a
			break
		} else {
			found = false
		}
	}

	if !found {
		log.Warnf("Action %s not found", step.ActionType)

		step.ActionMessages = append(step.ActionMessages, "Action not found")
		step.Running = false
		step.Error = true
		step.Finished = true
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(execution.ID.String(), step); err != nil {
			log.Error(err)
			return nil, false, false, false, false, err
		}

		return nil, false, false, false, true, nil
	}

	var flow_action models.Actions
	if len(flow.Actions) > 0 {
		for _, flowAction := range flow.Actions {
			if flowAction.ID.String() == step.ActionID {
				flow_action = flowAction
			}
		}
	}

	log.Info("Flow Action: ", flow_action)
	log.Infof("Execute Action: %s", action.Name)

	// resp, err := pluginManager.ExecutePlugin(action.ID, protocol.Request{
	// 	Action: "execute",
	// 	Data: map[string]interface{}{
	// 		"execution": execution,
	// 		"flow":      flow,
	// 		"payload":   payload,
	// 		"steps":     steps,
	// 		"step":      step,
	// 		"action":    flow_action,
	// 	},
	// })
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, false, false, false, true, err
	// }

	// if fn, ok := action.Function.(func(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool)); ok {
	// 	data, finished, canceled, no_pattern_match, failed := fn(execution, flow, payload, steps, step, flow_action)

	// 	if failed {
	// 		return nil, false, false, false, true, nil
	// 	} else if canceled {
	// 		return nil, false, true, false, false, nil
	// 	} else if no_pattern_match {
	// 		return nil, false, false, true, false, nil
	// 	} else if finished {
	// 		return data, true, false, false, false, nil
	// 	}
	// }

	// return data, true, false, false, false, nil

	return nil, false, false, false, false, nil
}
