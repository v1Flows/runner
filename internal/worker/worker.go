package worker

import (
	internal_executions "github.com/v1Flows/runner/internal/executions"
	"github.com/v1Flows/runner/pkg/plugins"
	shared_models "github.com/v1Flows/shared-library/pkg/models"
)

func StartWorker(platform string, actions []shared_models.Action, loadedPlugins map[string]plugins.Plugin) {
	internal_executions.GetPendingExecutions(platform, actions, loadedPlugins)
}
