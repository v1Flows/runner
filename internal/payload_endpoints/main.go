package payloadendpoints

import (
	"strconv"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/plugins"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func RegisterEndpoints(loadedPluginEndpoints []models.Plugins) (endpoints []models.PayloadEndpoints) {
	for _, plugin := range loadedPluginEndpoints {
		endpoints = append(endpoints, plugin.Endpoints)
	}

	if len(endpoints) == 0 {
		endpoints = []models.PayloadEndpoints{}
	}

	return endpoints
}

func InitPayloadRouter(cfg config.Config, endpointPlugins []models.Plugins, loadedPlugins map[string]plugins.Plugin) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Payload Port: ", cfg.PayloadEndpoints.Port)

	payload := router.Group("/payloads")
	for _, plugin := range endpointPlugins {
		log.Infof("Open %s Endpoint: %s", plugin.Name, plugin.Endpoints.Endpoint)
		payload.POST(plugin.Endpoints.Endpoint, func(c *gin.Context) {
			log.Info("Received Payload for: ", plugin.Name)
			loadedPlugins[plugin.Name].HandlePayload(plugins.PayloadHandlerRequest{
				Config:  cfg,
				Context: c,
			})
		})
	}

	router.Run(":" + strconv.Itoa(cfg.PayloadEndpoints.Port))
}
