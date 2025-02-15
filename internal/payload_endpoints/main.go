package payloadendpoints

import (
	"strconv"

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

func InitPayloadRouter(port int, endpointPlugins []models.Plugins) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Payload Port: ", port)

	payload := router.Group("/payloads")
	for _, plugin := range endpointPlugins {
		log.Infof("Open %s Endpoint: %s", plugin.Name, plugin.Endpoints.Endpoint)
		payload.POST(plugin.Endpoints.Endpoint, func(c *gin.Context) {
			log.Info("Received Payload: ", plugin.Name)
		})
	}

	router.Run(":" + strconv.Itoa(port))
}
