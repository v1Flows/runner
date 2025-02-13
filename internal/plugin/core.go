package plugin

import (
	"github.com/AlertFlow/runner/pkg/models"
	"github.com/gin-gonic/gin"
)

// Action defines the core interface that all plugins must implement
type Plugin interface {
	// Info returns metadata about the action
	Details() models.Plugin

	// Execute runs the action with given input and returns output or error
	Execute(c *gin.Context) error

	// Handles payload endpoints for incoming data
	Handle(c *gin.Context) error
}
