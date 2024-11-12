package payloadendpoints

import (
	"strconv"

	"gitlab.justlab.xyz/alertflow-public/runner/internal/plugins"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitPayloadRouter(port int, plugins []plugins.Plugin, payloadEndpoints []models.PayloadEndpoint) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Payload Port: ", port)

	payload := router.Group("/payloads")
	for _, endpoint := range payloadEndpoints {
		log.Infof("Open %s Endpoint: %s", endpoint.Name, endpoint.Endpoint)
		payload.POST(endpoint.Endpoint, func(c *gin.Context) {
			for _, plugin := range plugins {
				if plugin.Init().Name == endpoint.Name {
					plugin.Handle(c)
				}
			}
		})
	}

	router.Run(":" + strconv.Itoa(port))
}
