package payloadhandler

import (
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func InitPayloadRouter(port int, types []string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Payload Port: ", port)

	payload := router.Group("/payloads")
	{
		if slices.Contains(types, "Alertmanager") {
			log.Info("Open Alertmanager Endpoint: /alertmanager")
			payload.POST("/alertmanager", func(c *gin.Context) {
				AlertmanagerPayloadHandler(c)
			})
		}
	}

	router.Run(":" + strconv.Itoa(port))
}
