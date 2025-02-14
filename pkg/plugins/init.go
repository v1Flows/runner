// filepath: /Users/Justin.Neubert/projects/v1flows/alertflow/runner/pkg/plugins/init.go
package plugins

import (
	"log"
	"os/exec"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"
	"github.com/hashicorp/go-plugin"
)

var loadedPlugins = make(map[string]Plugin)
var actionPlugins = make([]models.Plugin, 0)
var endpointPlugins = make([]models.Plugin, 0)
var pluginClients = make([]*plugin.Client, 0) // Track plugin clients

func Init(cfg config.Config) (loadedPlugin map[string]Plugin, actionPlugins []models.Plugin, endpointPlugins []models.Plugin) {
	pluginPaths, err := DownloadAndBuildPlugins(cfg.Plugins, ".plugins_temp", cfg.PluginDir)
	if err != nil {
		log.Fatalf("Error downloading and building plugins: %v", err)
	}

	for name, path := range pluginPaths {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "PLUGIN_MAGIC_COOKIE",
				MagicCookieValue: "hello",
			},
			Plugins: map[string]plugin.Plugin{
				"plugin": &PluginServer{},
			},
			Cmd: exec.Command(path),
		})

		rpcClient, err := client.Client()
		if err != nil {
			log.Fatalf("Error loading plugin %s: %v", name, err)
		}

		raw, err := rpcClient.Dispense("plugin")
		if err != nil {
			log.Fatalf("Error dispensing plugin %s: %v", name, err)
		}

		plugin := raw.(Plugin)
		loadedPlugins[name] = plugin
		pluginClients = append(pluginClients, client) // Store the client

		// Get plugin info
		info, err := plugin.Info()
		if err != nil {
			log.Fatalf("Error getting info for plugin %s: %v", name, err)
		}

		if info.Type == "payload_endpoint" {
			endpointPlugins = append(endpointPlugins, info)
		} else if info.Type == "action" {
			actionPlugins = append(actionPlugins, info)
		}
	}

	return loadedPlugins, actionPlugins, endpointPlugins
}

// ShutdownPlugins terminates all plugin clients
func ShutdownPlugins() {
	for _, client := range pluginClients {
		client.Kill()
	}
}
