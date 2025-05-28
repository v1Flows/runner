package internal_executions

import (
	"errors"
	"net"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
	af_models "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/executions"
	"github.com/v1Flows/runner/pkg/platform"
	"github.com/v1Flows/runner/pkg/plugins"
	"github.com/v1Flows/shared-library/pkg/models"
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

func processStep(
	cfg config.Config,
	workspace string,
	actions []shared_models.Action,
	loadedPlugins map[string]plugins.Plugin,
	flow shared_models.Flows,
	flowBytes []byte,
	alert af_models.Alerts,
	steps []shared_models.ExecutionSteps,
	step shared_models.ExecutionSteps,
	execution shared_models.Executions,
	broker *plugin.MuxBroker,
) (res plugins.Response, success bool, canceled bool, err error) {
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
		return plugins.Response{}, false, false, err
	}

	valid, danger, pluginVersion := common.CheckActionVersionAgainstPluginVersion(actions, step)

	if !valid {
		// dont execute step and quit execution
		step.Messages = append(step.Messages, shared_models.Message{
			Title: "Error",
			Lines: []shared_models.Line{
				{
					Content:   "Action not compatible with plugin version",
					Color:     "danger",
					Timestamp: time.Now(),
				},
				{
					Content:   "Plugin Version: " + pluginVersion,
					Color:     "danger",
					Timestamp: time.Now(),
				},
				{
					Content:   "Action Version: " + step.Action.Version,
					Color:     "danger",
					Timestamp: time.Now(),
				},
				{
					Content:   "Cancel execution",
					Color:     "danger",
					Timestamp: time.Now(),
				},
			},
		})
		step.Status = "error"
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(cfg, execution.ID.String(), step, targetPlatform); err != nil {
			log.Error(err)
			return plugins.Response{}, false, false, err
		}

		return plugins.Response{}, false, false, nil
	}

	if danger {
		// modify the pickup step
		err = executions.UpdateStep(cfg, execution.ID.String(), models.ExecutionSteps{
			ID: steps[0].ID,
			Messages: []models.Message{
				{
					Title: "Caution",
					Lines: []models.Line{
						{
							Content:   "Plugin version is higher than action version. This may cause issues but execution will still be processed.",
							Timestamp: time.Now(),
							Color:     "warning",
						},
					},
				},
			},
			Status: "running",
		}, targetPlatform)
		if err != nil {
			return plugins.Response{}, false, false, err
		}
	}

	if _, ok := loadedPlugins[step.Action.Plugin]; !ok {
		log.Warnf("Action %s not found", step.Action.Plugin)

		step.Messages = append(step.Messages, shared_models.Message{
			Title: "Error",
			Lines: []shared_models.Line{
				{
					Content:   "Action not found in loaded plugins",
					Color:     "danger",
					Timestamp: time.Now(),
				},
				{
					Content:   "Target plugin: " + step.Action.Plugin,
					Color:     "danger",
					Timestamp: time.Now(),
				},
				{
					Content:   "Cancel execution",
					Color:     "danger",
					Timestamp: time.Now(),
				},
			},
		})
		step.Status = "error"
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(cfg, execution.ID.String(), step, targetPlatform); err != nil {
			log.Error(err)
			return plugins.Response{}, false, false, err
		}

		return plugins.Response{}, false, false, errors.New("plugin not found")
	}

	// Register the runner RPC server and get a brokerID
	brokerID := broker.NextId()
	go broker.AcceptAndServe(brokerID, func(conn net.Conn) {
		rpc.Register(&plugins.RunnerRPCServer{})
		rpc.ServeConn(conn)
	})

	req := plugins.ExecuteTaskRequest{
		Flow:      flow,
		FlowBytes: flowBytes,
		Execution: execution,
		Step:      step,
		Alert:     alert,
		Platform:  targetPlatform,
		Workspace: workspace,
		BrokerID:  int(brokerID),
	}

	res, err = loadedPlugins[step.Action.Plugin].ExecuteTask(req)
	if err != nil {
		log.Error(err)

		step.Messages = append(step.Messages, shared_models.Message{
			Title: "Error",
			Lines: []shared_models.Line{
				{
					Content:   "Failed to execute action",
					Color:     "danger",
					Timestamp: time.Now(),
				},
				{
					Content:   "Error: " + err.Error(),
					Color:     "danger",
					Timestamp: time.Now(),
				},
				{
					Content:   "Cancel execution",
					Color:     "danger",
					Timestamp: time.Now(),
				},
			},
		})
		step.Status = "error"
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(cfg, execution.ID.String(), step, targetPlatform); err != nil {
			log.Error(err)
			return plugins.Response{}, false, false, err
		}

		return plugins.Response{}, false, false, err
	}

	if res.Canceled {
		return res, false, true, nil
	} else if res.Success {
		return res, true, false, nil
	} else {
		return res, false, false, nil
	}

	// return data, true, false, false, false, nil
}
