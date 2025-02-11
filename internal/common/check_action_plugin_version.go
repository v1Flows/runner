package common

import (
	"github.com/AlertFlow/runner/internal/plugins"
	"github.com/AlertFlow/runner/pkg/models"
)

func checkActionVersionAgainstPluginVersion(step models.ExecutionSteps) (valid bool, pluginVersion string) {
	pluginVersion = plugins.GetPluginVersion(step.ActionType)

	// Remove the 'v' prefix from the plugin version if it exists
	if len(pluginVersion) > 0 && pluginVersion[0] == 'v' {
		pluginVersion = pluginVersion[1:]
	}

	return pluginVersion == step.ActionVersion, pluginVersion
}
