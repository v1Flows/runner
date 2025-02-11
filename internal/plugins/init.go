package plugins

import (
	"os"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func Init() ([]Plugin, []models.Plugin, []models.ActionDetails, []models.PayloadEndpoint) {
	pluginDir := "plugins"
	pluginTempDir := "plugins_temp"

	for _, plugin := range config.Config.Plugins {
		// Check if the plugin is already installed and up-to-date
		if isPluginPresent(plugin.Name, pluginDir) && isPluginUpToDate(plugin.Name, plugin.Version, pluginDir) {
			log.Infof("Plugin %s is already up-to-date", plugin.Name)
			continue
		} else {
			// Remove the plugin from the versions file
			err := removePluginFromVersionsFile(plugin.Name, pluginDir)
			if err != nil {
				log.Errorf("Failed to remove plugin %s from versions file: %v", plugin.Name, err)
			}
		}

		// Clone and build the plugin
		log.Infof("Cloning and building plugin %s", plugin.Name)
		err := cloneAndBuildPlugin(plugin.Url, pluginDir, pluginTempDir, plugin.Name, plugin.Version)
		if err != nil {
			log.Errorf("Failed to clone and build plugin %s: %v", plugin.Name, err)
		}
	}

	// cleanup the temp directory
	err := os.RemoveAll(pluginTempDir)
	if err != nil {
		log.Errorf("Failed to remove temp directory: %v", err)
	}

	plugins, err := loadPlugins(pluginDir)
	if err != nil {
		log.Fatal(err)
	}

	pluginsMap := []models.Plugin{}
	actions := make([]models.ActionDetails, 0)
	payloadEndpoints := make([]models.PayloadEndpoint, 0)
	for _, plugin := range plugins {
		p := plugin.Init()

		pluginsMap = append(pluginsMap, p)

		if p.Type == "action" {
			action := plugin.Details()
			action.Action.Version = p.Version
			actions = append(actions, action.Action)
			log.Infof("Loaded action plugin: %s", action.Action.Name)
		}
		if p.Type == "payload_endpoint" {
			payloadEndpoint := plugin.Details()
			payloadEndpoints = append(payloadEndpoints, payloadEndpoint.Payload)
			log.Infof("Loaded payload endpoint plugin: %s", payloadEndpoint.Payload.Name)
		}
	}

	return plugins, pluginsMap, actions, payloadEndpoints
}
