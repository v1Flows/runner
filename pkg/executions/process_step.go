package executions

import (
	"errors"
	"time"

	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/plugins"

	log "github.com/sirupsen/logrus"
)

func RegisterActions(loadedPluginActions []models.Plugins) (actions []models.Actions) {
	for _, plugin := range loadedPluginActions {
		actions = append(actions, plugin.Actions)
	}

	if len(actions) == 0 {
		actions = []models.Actions{}
	}

	return actions
}

func processStep(cfg config.Config, actions []models.Actions, loadedPlugins map[string]plugins.Plugin, flow models.Flows, alert models.Alerts, steps []models.ExecutionSteps, step models.ExecutionSteps, execution models.Executions) (res plugins.Response, success bool, err error) {
	step.Status = "running"
	step.StartedAt = time.Now()
	step.RunnerID = execution.RunnerID

	if err := UpdateStep(cfg, execution.ID.String(), step); err != nil {
		log.Error(err)
		return plugins.Response{}, false, err
	}

	valid, pluginVersion := common.CheckActionVersionAgainstPluginVersion(actions, step)

	if !valid {
		// dont execute step and quit execution
		step.Messages = append(step.Messages, "Action not compatible with plugin version", "Plugin Version: "+pluginVersion+" Action Version: "+step.Action.Version, "Stopping execution")
		step.Status = "error"
		step.FinishedAt = time.Now()

		if err := UpdateStep(cfg, execution.ID.String(), step); err != nil {
			log.Error(err)
			return plugins.Response{}, false, err
		}

		return plugins.Response{}, false, nil
	}

	if _, ok := loadedPlugins[step.Action.Plugin]; !ok {
		log.Warnf("Action %s not found", step.Action.Plugin)

		step.Messages = append(step.Messages, "Action not found in loaded plugins", "Target plugin: "+step.Action.Plugin, "Stopping execution")
		step.Status = "error"
		step.FinishedAt = time.Now()

		if err := UpdateStep(cfg, execution.ID.String(), step); err != nil {
			log.Error(err)
			return plugins.Response{}, false, err
		}

		return plugins.Response{}, false, errors.New("plugin not found")
	}

	req := plugins.ExecuteTaskRequest{
		Config:    cfg,
		Flow:      flow,
		Execution: execution,
		Step:      step,
		Alert:     alert,
	}

	res, err = loadedPlugins[step.Action.Plugin].ExecuteTask(req)
	if err != nil {
		log.Error(err)
		return plugins.Response{}, false, err
	}

	if res.Success {
		return res, true, nil
	} else {
		return res, false, nil
	}

	// return data, true, false, false, false, nil
}
