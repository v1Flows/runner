package payloadhandler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func InitPayloadRouter(port int, types []string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Payload Port: ", port)

	// payload := router.Group("/payloads")
	// {
	// 	if slices.Contains(types, "Alertmanager") {
	// 		log.Info("Open Alertmanager Endpoint: /alertmanager")
	// 	}
	// }

	router.Run(":" + strconv.Itoa(port))
}
