package plugin

import "github.com/AlertFlow/runner/config"

func GetPluginVersion(plugin string) string {
	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	// find plugin in config.Config.Plugins
	// return plugin.Version
	for _, p := range cfg.Plugins {
		if p.Name == plugin {
			return p.Version
		}
	}

	return ""
}
