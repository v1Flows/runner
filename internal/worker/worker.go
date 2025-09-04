package worker

import (
	"github.com/v1Flows/exFlow/services/backend/pkg/models"
	internal_executions "github.com/v1Flows/runner/internal/executions"
	"github.com/v1Flows/runner/pkg/plugins"
)

func StartWorker(actions []models.Action, loadedPlugins map[string]plugins.Plugin) {
	internal_executions.GetPendingExecutions(actions, loadedPlugins)
}
