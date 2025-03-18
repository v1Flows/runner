package endpoints

import (
	"fmt"
	"io"
	"strconv"

	af_models "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/flows"
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

func InitEndpointRouter(cfg config.Config, router *gin.Engine, platform string, endpointPlugins []shared_models.Plugin, loadedPlugins map[string]plugins.Plugin) {
	log.Info("Open Alert Port: ", cfg.Endpoints.Port)

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

			var afFlow af_models.Flows
			_, afFlow, err = flows.GetFlowData(cfg, "50c9ce52-0590-47a6-a4f4-cb7d3b632b47", platform)
			if err != nil {
				log.Error("Error getting flow data: ", err)
			}
			fmt.Println("AlertFlow Flow: ", afFlow.GroupAlerts)

			request := plugins.EndpointRequest{
				Config:   cfg,
				Body:     bodyBytes,
				Platform: platform,
			}

			res, err := loadedPlugins[plugin.Endpoint.ID].EndpointRequest(request)
			if err != nil {
				log.Error("Error in handling request: ", err)
				c.JSON(500, gin.H{
					"error": err,
				})
			} else {
				log.Info("Request handled successfully")
				c.JSON(200, gin.H{
					"response": res,
				})
			}
		})
	}

	router.Run(":" + strconv.Itoa(cfg.Endpoints.Port))
}
