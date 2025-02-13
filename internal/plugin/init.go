package plugin

import (
	"context"
	"fmt"
	"log"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"
)

func Init(config config.Config) (manager *Manager, plugins []models.Plugin, actions []models.ActionDetails, payloadEndpoints []models.PayloadEndpoint) {
	manager, err := NewManager(config)
	if err != nil {
		log.Fatal(err)
	}

	actions = make([]models.ActionDetails, 0)
	payloadEndpoints = make([]models.PayloadEndpoint, 0)

	for _, plugin := range config.Plugins {
		if err := manager.InstallPlugin(context.Background(), plugin.Name); err != nil {
			log.Fatal(err)
		}

		client, err := manager.LoadPlugin(plugin.Name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(client)
		fmt.Println(manager)

		// plugins = append(plugins, plugin.Name)

		// if resp.Plugin.Action is defined with content, add it to actions
		// if resp.Plugin.Action.ID != "" {
		// 	resp.Plugin.Action.Version = resp.Plugin.Version
		// 	actions = append(actions, resp.Plugin.Action)
		// }

		// if resp.Plugin.Payload.Endpoint != "" {
		// 	resp.Plugin.Payload.Version = resp.Plugin.Version
		// 	payloadEndpoints = append(payloadEndpoints, resp.Plugin.Payload)
		// }
	}

	return manager, plugins, actions, payloadEndpoints
}
