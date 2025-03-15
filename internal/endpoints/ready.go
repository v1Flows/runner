package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/v1Flows/runner/config"
)

func ReadyEndpoint(cfg config.Config, router *gin.Engine) {
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
