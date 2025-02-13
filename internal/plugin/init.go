package plugin

import (
	"github.com/AlertFlow/runner/pkg/models"
)

func Init() ([]models.Plugin, []models.ActionDetails, []models.PayloadEndpoint) {
	// pluginDir := "plugins"
	// pluginTempDir := "plugins_temp"

	registry := newRegistry()

	logPlugin := logPlugin.New()
	registry.Register(&ActionPlugin{})

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
