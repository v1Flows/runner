// filepath: /Users/Justin.Neubert/projects/v1flows/alertflow/runner/pkg/plugins/init.go
package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/AlertFlow/runner/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

var loadedPlugins = make(map[string]Plugin)
var pluginClients = make([]*plugin.Client, 0) // Track plugin clients

const maxRetries = 3
const retryInterval = 5 * time.Second

func connectPlugin(name, path string) (Plugin, *plugin.Client, error) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   fmt.Sprintf("plugin.%s", name),
		Output: os.Stdout,
		Level:  hclog.Error, // Set to Error level to suppress debug logs
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "PLUGIN_MAGIC_COOKIE",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			"plugin": &PluginServer{},
		},
		Cmd:    exec.Command(path),
		Logger: logger, // Suppress plugin logs
	})

	var rpcClient plugin.ClientProtocol
	var err error
	for i := 0; i < maxRetries; i++ {
		rpcClient, err = client.Client()
		if err == nil {
			break
		}
		log.Errorf("Error loading plugin %s: %v. Retrying in %v...", name, err, retryInterval)
		time.Sleep(retryInterval)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to plugin %s after %d retries: %v", name, maxRetries, err)
	}

	raw, err := rpcClient.Dispense("plugin")
	if err != nil {
		return nil, nil, fmt.Errorf("error dispensing plugin %s: %v", name, err)
	}

	plugin := raw.(Plugin)
	return plugin, client, nil
}

func Init(cfg config.Config) (loadedPlugin map[string]Plugin, plugins []models.Plugins, actionPlugins []models.Plugins, endpointPlugins []models.Plugins) {
	// Define mandatory plugins
	mandatoryPlugins := []config.PluginConfig{
		{Name: "collect_data", Version: "main", Repository: "https://github.com/AlertFlow/rp-collect_data"},
		{Name: "pattern_check", Version: "main", Repository: "https://github.com/AlertFlow/rp-pattern_check"},
		{Name: "actions_check", Version: "main", Repository: "https://github.com/AlertFlow/rp-actions_check"},
		{Name: "log", Version: "main", Repository: "https://github.com/AlertFlow/rp-log"},
		{Name: "wait", Version: "main", Repository: "https://github.com/AlertFlow/rp-wait"},
		{Name: "interaction", Version: "main", Repository: "https://github.com/AlertFlow/rp-interaction"},
		{Name: "ping", Version: "main", Repository: "https://github.com/AlertFlow/rp-ping"},
		{Name: "port_checker", Version: "main", Repository: "https://github.com/AlertFlow/rp-port_checker"},
	}

	// Merge mandatory plugins with config plugins, handling conflicts
	pluginMap := make(map[string]config.PluginConfig)
	for _, plugin := range mandatoryPlugins {
		pluginMap[plugin.Name] = plugin
	}
	for _, plugin := range cfg.Plugins {
		if _, exists := pluginMap[plugin.Name]; exists {
			log.Warnf("Conflict: Plugin %s is defined in both mandatory list and config. Using config version.", plugin.Name)
		}
		pluginMap[plugin.Name] = plugin
	}

	// Convert pluginMap to a slice
	var allPlugins []config.PluginConfig
	for _, plugin := range pluginMap {
		allPlugins = append(allPlugins, plugin)
	}

	pluginPaths, err := DownloadAndBuildPlugins(allPlugins, ".plugins_temp", cfg.PluginDir)
	if err != nil {
		log.Fatalf("Error downloading and building plugins: %v", err)
	}

	err = CleanupUnusedPlugins(allPlugins, cfg.PluginDir)
	if err != nil {
		log.Warnf("Error cleaning up unused plugins: %v", err)
	}

	for name, path := range pluginPaths {
		plugin, client, err := connectPlugin(name, path)
		if err != nil {
			log.Fatalf("Error connecting to plugin %s: %v", name, err)
		}

		loadedPlugins[name] = plugin
		pluginClients = append(pluginClients, client) // Store the client

		// Get plugin info
		info, err := plugin.Info()
		if err != nil {
			log.Fatalf("Error getting info for plugin %s: %v", name, err)
		}

		plugins = append(plugins, info)
		if info.Type == "endpoint" {
			endpointPlugins = append(endpointPlugins, info)
		} else if info.Type == "action" {

			if info.Actions.Version == "" {
				info.Actions.Version = info.Version
			}

			actionPlugins = append(actionPlugins, info)
		}
	}

	return loadedPlugins, plugins, actionPlugins, endpointPlugins
}

// ShutdownPlugins terminates all plugin clients
func ShutdownPlugins() {
	for _, client := range pluginClients {
		client.Kill()
	}
}
