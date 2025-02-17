package common

import (
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

func checkActionVersionAgainstPluginVersion(actions []models.Actions, step bmodels.ExecutionSteps) (valid bool, pluginVersion string) {
	for _, action := range actions {
		if action.Plugin == step.Action.Plugin {
			pluginVersion = action.Version
			break
		}
	}

	// Remove the 'v' prefix from the plugin version if it exists
	if len(pluginVersion) > 0 && pluginVersion[0] == 'v' {
		pluginVersion = pluginVersion[1:]
	}

	if step.Action.Version == "" {
		return true, pluginVersion
	}

	return pluginVersion == step.Action.Version, pluginVersion
}
