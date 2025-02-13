package payloadendpoints

import (
	"strconv"

	"github.com/AlertFlow/runner/internal/plugin"
	"github.com/AlertFlow/runner/pkg/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitPayloadRouter(port int, plugins plugin.Plugin, payloadEndpoints []models.PayloadEndpoint) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Payload Port: ", port)

	payload := router.Group("/payloads")
	for _, endpoint := range payloadEndpoints {
		log.Infof("Open %s Endpoint: %s", endpoint.Name, endpoint.Endpoint)
		payload.POST(endpoint.Endpoint, func(c *gin.Context) {
			if plugin.Details().Name == endpoint.Name {
				plugin.Execute(c)
			}
		})
	}

	router.Run(":" + strconv.Itoa(port))
}
