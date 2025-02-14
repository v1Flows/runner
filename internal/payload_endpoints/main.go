package payloadendpoints

import (
	"strconv"

	"github.com/AlertFlow/runner/pkg/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitPayloadRouter(port int, endpointPlugins []models.Plugin) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Payload Port: ", port)

	payload := router.Group("/payloads")
	for _, endpoint := range endpointPlugins {
		log.Infof("Open %s Endpoint: %s", endpoint.Name, endpoint.Payload.Endpoint)
		payload.POST(endpoint.Payload.Endpoint, func(c *gin.Context) {
			log.Info("Received Payload: ", endpoint.Name)
		})
	}

	router.Run(":" + strconv.Itoa(port))
}
