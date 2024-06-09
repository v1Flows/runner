package incoming

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func InitPayloadRouter(port int) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Payload Port: ", port)

	payload := router.Group("/payload")
	{
		payload.POST("/alertmanager", func(c *gin.Context) {
			AlertmanagerPayloadHandler(c)
		})

		payload.POST("/zabbix", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "alive",
			})
		})
	}

	router.Run(":" + strconv.Itoa(port))
}
