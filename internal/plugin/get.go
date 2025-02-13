package plugin

import "github.com/AlertFlow/runner/config"

func GetPluginVersion(plugin string) string {
	// find plugin in config.Config.Plugins
	// return plugin.Version
	for _, p := range config.Config.Plugins {
		if p.Name == plugin {
			return p.Version
		}
	}

	return ""
}
