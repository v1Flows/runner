package internal_executions

import (
	"errors"
	"time"

	af_models "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/executions"
	"github.com/v1Flows/runner/pkg/platform"
	"github.com/v1Flows/runner/pkg/plugins"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func RegisterActions(loadedPluginActions []shared_models.Plugin) (actions []shared_models.Action) {
	for _, plugin := range loadedPluginActions {
		actions = append(actions, plugin.Action)
	}

	if len(actions) == 0 {
		actions = []shared_models.Action{}
	}

	return actions
}

func processStep(cfg config.Config, workspace string, actions []shared_models.Action, loadedPlugins map[string]plugins.Plugin, flow shared_models.Flows, flowBytes []byte, alert af_models.Alerts, steps []shared_models.ExecutionSteps, step shared_models.ExecutionSteps, execution shared_models.Executions) (res plugins.Response, success bool, err error) {
	targetPlatform, ok := platform.GetPlatformForExecution(execution.ID.String())
	if !ok {
		log.Error("Failed to get platform")
		return
	}

	step.Status = "running"
	step.StartedAt = time.Now()
	step.RunnerID = execution.RunnerID

	if err := executions.UpdateStep(cfg, execution.ID.String(), step, targetPlatform); err != nil {
		log.Error(err)
		return plugins.Response{}, false, err
	}

	valid, pluginVersion := common.CheckActionVersionAgainstPluginVersion(actions, step)

	if !valid {
		// dont execute step and quit execution
		step.Messages = append(step.Messages, shared_models.Message{
			Title: "Error",
			Lines: []shared_models.Line{
				{
					Content: "Action not compatible with plugin version",
					Color:   "danger",
				},
				{
					Content: "Plugin Version: " + pluginVersion,
					Color:   "danger",
				},
				{
					Content: "Action Version: " + step.Action.Version,
					Color:   "danger",
				},
				{
					Content: "Cancel execution",
					Color:   "danger",
				},
			},
		})
		step.Status = "error"
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(cfg, execution.ID.String(), step, targetPlatform); err != nil {
			log.Error(err)
			return plugins.Response{}, false, err
		}

		return plugins.Response{}, false, nil
	}

	if _, ok := loadedPlugins[step.Action.Plugin]; !ok {
		log.Warnf("Action %s not found", step.Action.Plugin)

		step.Messages = append(step.Messages, shared_models.Message{
			Title: "Error",
			Lines: []shared_models.Line{
				{
					Content: "Action not found in loaded plugins",
					Color:   "danger",
				},
				{
					Content: "Target plugin: " + step.Action.Plugin,
					Color:   "danger",
				},
				{
					Content: "Cancel execution",
					Color:   "danger",
				},
			},
		})
		step.Status = "error"
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(cfg, execution.ID.String(), step, targetPlatform); err != nil {
			log.Error(err)
			return plugins.Response{}, false, err
		}

		return plugins.Response{}, false, errors.New("plugin not found")
	}

	req := plugins.ExecuteTaskRequest{
		Config:    cfg,
		Flow:      flow,
		FlowBytes: flowBytes,
		Execution: execution,
		Step:      step,
		Alert:     alert,
		Platform:  targetPlatform,
		Workspace: workspace,
	}

	res, err = loadedPlugins[step.Action.Plugin].ExecuteTask(req)
	if err != nil {
		log.Error(err)

		step.Messages = append(step.Messages, shared_models.Message{
			Title: "Error",
			Lines: []shared_models.Line{
				{
					Content: "Failed to execute action",
					Color:   "danger",
				},
				{
					Content: "Error: " + err.Error(),
					Color:   "danger",
				},
				{
					Content: "Cancel execution",
					Color:   "danger",
				},
			},
		})
		step.Status = "error"
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(cfg, execution.ID.String(), step, targetPlatform); err != nil {
			log.Error(err)
			return plugins.Response{}, false, err
		}

		return plugins.Response{}, false, err
	}

	if res.Success {
		return res, true, nil
	} else {
		return res, false, nil
	}

	// return data, true, false, false, false, nil
}
