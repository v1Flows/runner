package plugins

import (
	"path/filepath"
	"plugin"

	"github.com/AlertFlow/runner/pkg/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Plugin interface {
	Init() models.Plugin
	Details() models.PluginDetails
	Execute(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool)
	Handle(context *gin.Context)
}

func loadPlugins(pluginDir string) ([]Plugin, error) {
	var plugins []Plugin

	log.Info("Loading plugins from ", pluginDir)

	files, err := filepath.Glob(filepath.Join(pluginDir, "*.so"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p, err := plugin.Open(file)
		if err != nil {
			log.Errorln("Error loading plugin:", err)
			continue
		}

		sym, err := p.Lookup("Plugin")
		if err != nil {
			log.Errorln("Error looking up symbol:", err)
			continue
		}

		Plugin, ok := sym.(Plugin)
		if !ok {
			log.Errorln("Invalid plugin type in", file)
			continue
		}

		plugins = append(plugins, Plugin)
	}

	return plugins, nil
}
