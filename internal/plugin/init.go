package plugin

import (
	"fmt"
	"log"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"
	"github.com/AlertFlow/runner/pkg/protocol"
)

func Init() ([]models.Plugin, []models.ActionDetails, []models.PayloadEndpoint) {
	// pluginDir := "plugins"
	// pluginTempDir := "plugins_temp"

	manager := NewManager("plugins", "plugins_temp")

	if err := manager.DownloadPlugin(config.PluginConf{
		Name:    "Log",
		Url:     "https://github.com/AlertFlow/rp-log",
		Version: "refactor",
	}); err != nil {
		log.Fatal(err)
	}

	if err := manager.StartPlugin(config.PluginConf{Name: "Log"}); err != nil {
		log.Fatal(err)
	}

	// Execute plugin
	resp, err := manager.ExecutePlugin("Log", protocol.Request{
		Action: "details",
		Data: map[string]interface{}{
			"param1": "value1",
		},
	})

	fmt.Println(resp.Plugin, err)

	resp, err = manager.ExecutePlugin("Log", protocol.Request{
		Action: "process",
		Data: map[string]interface{}{
			"param1": "value1",
		},
	})

	fmt.Println(resp, err)

	// pluginsMap := []models.Plugin{}
	// actions := make([]models.ActionDetails, 0)
	// payloadEndpoints := make([]models.PayloadEndpoint, 0)
	// for _, plugin := range plugins {
	// 	p := plugin.Init()

	// 	pluginsMap = append(pluginsMap, p)

	// 	if p.Type == "action" {
	// 		action := plugin.Details()
	// 		action.Action.Version = p.Version
	// 		actions = append(actions, action.Action)
	// 		log.Infof("Loaded action plugin: %s", action.Action.Name)
	// 	}
	// 	if p.Type == "payload_endpoint" {
	// 		payloadEndpoint := plugin.Details()
	// 		payloadEndpoints = append(payloadEndpoints, payloadEndpoint.Payload)
	// 		log.Infof("Loaded payload endpoint plugin: %s", payloadEndpoint.Payload.Name)
	// 	}
	// }

	return nil, nil, nil
}
