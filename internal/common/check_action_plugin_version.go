package common

import (
	"github.com/AlertFlow/runner/pkg/models"
	"github.com/AlertFlow/runner/pkg/plugin"
)

func checkActionVersionAgainstPluginVersion(step models.ExecutionSteps) (valid bool, pluginVersion string) {
	pluginVersion = plugin.GetPluginVersion(step.ActionType)

	// Remove the 'v' prefix from the plugin version if it exists
	if len(pluginVersion) > 0 && pluginVersion[0] == 'v' {
		pluginVersion = pluginVersion[1:]
	}

	if step.ActionVersion == "" {
		return true, pluginVersion
	}

	return pluginVersion == step.ActionVersion, pluginVersion
}
