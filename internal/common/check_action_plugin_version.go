package common

import (
	shared_models "github.com/v1Flows/shared-library/pkg/models"
)

func CheckActionVersionAgainstPluginVersion(actions []shared_models.Action, step shared_models.ExecutionSteps) (valid bool, pluginVersion string) {
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

	if pluginVersion < step.Action.Version {
		return false, pluginVersion
	}

	return true, pluginVersion
}
