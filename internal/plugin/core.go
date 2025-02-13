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

	// Validate validates the plugin configuration
	Validate() error
}
