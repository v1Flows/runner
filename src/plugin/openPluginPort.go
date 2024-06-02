package plugin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func OpenPluginPort(port int) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	log.Info("Open Plugin Port: ", port)

	runner := router.Group("/runner")
	{
		runner.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "alive",
			})
		})

		publish := runner.Group("/publish")
		{
			publish.POST("/action", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "alive",
				})
			})
			publish.POST("/payload", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "alive",
				})
			})
		}

		execution := runner.Group("/execution")
		{
			execution.GET("/pending", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "alive",
				})
			})
			execution.POST("/finish", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "alive",
				})
			})
			execution.POST("/start", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "alive",
				})
			})
		}
	}

	router.Run(":" + strconv.Itoa(port))
}
