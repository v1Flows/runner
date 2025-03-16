package worker

import (
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/executions"
	"github.com/v1Flows/runner/pkg/plugins"
)

func StartWorker(platform string, cfg config.Config, actions []models.Actions, loadedPlugins map[string]plugins.Plugin) {
	executions.GetPendingExecutions(platform, cfg, actions, loadedPlugins)
}
