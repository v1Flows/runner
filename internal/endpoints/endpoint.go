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

func InitRouter(cfg config.Config, router *gin.Engine, platform string, endpointPlugins []shared_models.Plugin, loadedPlugins map[string]plugins.Plugin) {
	// enable api endpoint for backend -> runner communication
	log.Info("Router Listening on Port: ", cfg.ApiEndpoint.Port)
	v1 := router.Group("/api/v1")

	v1.POST("/execution/:executionID/cancel", func(c *gin.Context) {
		log.Info("Received Cancel Request for Execution ID: ", c.Param("executionID"))
	})

	// handle incoming alert requests for alertflow
	if platform == "alertflow" && (cfg.Mode == "listener" || cfg.Mode == "master") {
		log.Info("Open Alert Port: ", cfg.ApiEndpoint.Port)

		alert := v1.Group("/alert")
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
	}

	router.Run(":" + strconv.Itoa(cfg.ApiEndpoint.Port))
}
