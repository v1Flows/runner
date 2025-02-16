package common

import (
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

func checkActionVersionAgainstPluginVersion(step bmodels.ExecutionSteps) (valid bool, pluginVersion string) {
	pluginVersion = "1.0.0"

	// Remove the 'v' prefix from the plugin version if it exists
	if len(pluginVersion) > 0 && pluginVersion[0] == 'v' {
		pluginVersion = pluginVersion[1:]
	}

	if step.Action.Version == "" {
		return true, pluginVersion
	}

	return pluginVersion == step.Action.Version, pluginVersion
}
