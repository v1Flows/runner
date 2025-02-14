// filepath: /Users/Justin.Neubert/projects/v1flows/alertflow/runner/pkg/plugins/init.go
package plugins

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/AlertFlow/runner/config"
	"github.com/hashicorp/go-plugin"
)

var loadedPlugins = make(map[string]Plugin)

func Init(cfg config.Config) {

	pluginPaths, err := DownloadAndBuildPlugins(cfg.Plugins, ".plugins_temp", cfg.PluginDir)
	if err != nil {
		log.Fatalf("Error downloading and building plugins: %v", err)
	}

	fmt.Println("Plugin Paths: ", pluginPaths)

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

		defer client.Kill()

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

		// Execute the plugin
		result, err := plugin.Execute(map[string]string{"target": "example.com"})
		if err != nil {
			log.Fatalf("Error executing plugin %s: %v", name, err)
		}

		fmt.Printf("Plugin %s Execute Result: %s\n", name, result)

		// Get plugin info
		info, err := plugin.Info()
		if err != nil {
			log.Fatalf("Error getting info for plugin %s: %v", name, err)
		}

		fmt.Printf("Plugin %s Info: %+v\n", name, info)
	}
}

func GetLoadedPlugins() map[string]Plugin {
	return loadedPlugins
}
