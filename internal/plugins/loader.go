package plugins

import (
	"os"
	"path/filepath"
	"plugin"
	"strings"

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

	pluginRemove := false

	for _, file := range files {
		p, err := plugin.Open(file)
		if err != nil {
			log.Error("Error opening plugin:", err.Error())
			if strings.Contains(err.Error(), "plugin was built with a different version of package") {
				pluginRemove = true
				os.Remove(file) // Assuming pluginPath is the path to the plugin file
				log.Warnln("Removed plugin due to version mismatch:", file)
			} else {
				log.Errorln("Error loading plugin:", err)
			}
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

	if pluginRemove {
		Init()
	}

	return plugins, nil
}
