package worker

import (
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	internal_executions "github.com/v1Flows/runner/internal/executions"
	"github.com/v1Flows/runner/pkg/plugins"
)

func StartWorker(platform string, cfg config.Config, actions []models.Actions, loadedPlugins map[string]plugins.Plugin) {
	internal_executions.GetPendingExecutions(platform, cfg, actions, loadedPlugins)
}
