package endpoints

import (
	"io"
	"strconv"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/plugins"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func RegisterEndpoints(loadedPluginEndpoints []shared_models.Plugin) (endpoints []shared_models.Endpoint) {
	for _, plugin := range loadedPluginEndpoints {
		endpoints = append(endpoints, plugin.Endpoint)
	}

	if len(endpoints) == 0 {
		endpoints = []shared_models.Endpoint{}
	}

	return endpoints
}

func InitAlertRouter(cfg config.Config, router *gin.Engine, endpointPlugins []shared_models.Plugin, loadedPlugins map[string]plugins.Plugin) {
	log.Info("Open Alert Port: ", cfg.AlertEndpoints.Port)

	alert := router.Group("/alert")
	for _, plugin := range endpointPlugins {
		log.Infof("Open %s Endpoint at /alert%s", plugin.Name, plugin.Endpoint.Path)
		alert.POST(plugin.Endpoint.Path, func(c *gin.Context) {
			log.Info("Received Alert for: ", plugin.Name)

			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Error("Error reading request body: ", err)
				c.JSON(500, gin.H{
					"error": "Error reading request body",
				})
				return
			}

			request := plugins.AlertHandlerRequest{
				Config: cfg,
				Body:   bodyBytes,
			}

			res, err := loadedPlugins[plugin.Endpoint.ID].HandleAlert(request)
			if err != nil {
				log.Error("Error in handling alert: ", err)
				c.JSON(500, gin.H{
					"error": err,
				})
			} else {
				log.Info("Alert handled successfully")
				c.JSON(200, gin.H{
					"response": res,
				})
			}
		})
	}

	router.Run(":" + strconv.Itoa(cfg.AlertEndpoints.Port))
}
