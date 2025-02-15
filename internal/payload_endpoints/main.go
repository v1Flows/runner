package payloadendpoints

import (
	"io"
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

			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Error("Error reading request body: ", err)
				c.JSON(500, gin.H{
					"error": "Error reading request body",
				})
				return
			}

			request := plugins.PayloadHandlerRequest{
				Config: cfg,
				Body:   bodyBytes,
			}

			res, err := loadedPlugins[plugin.Endpoints.ID].HandlePayload(request)
			if err != nil {
				log.Error("Error in handling payload: ", err)
				c.JSON(500, gin.H{
					"error": err,
				})
			} else {
				log.Info("Payload handled successfully")
				c.JSON(200, gin.H{
					"response": res,
				})
			}
		})
	}

	router.Run(":" + strconv.Itoa(cfg.PayloadEndpoints.Port))
}
