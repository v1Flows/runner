package plugins

import (
	"path/filepath"
	"plugin"

	log "github.com/sirupsen/logrus"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"
)

type ActionPlugin interface {
	Init() models.ActionDetails
	Execute(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool)
}

func LoadPlugins(pluginDir string) ([]ActionPlugin, error) {
	var plugins []ActionPlugin

	log.Info("Loading plugins from ", pluginDir)

	files, err := filepath.Glob(filepath.Join(pluginDir, "*.so"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p, err := plugin.Open(file)
		if err != nil {
			log.Println("Error loading plugin:", err)
			continue
		}

		sym, err := p.Lookup("Plugin")
		if err != nil {
			log.Println("Error looking up symbol:", err)
			continue
		}

		actionPlugin, ok := sym.(ActionPlugin)
		if !ok {
			log.Println("Invalid plugin type")
			continue
		}

		plugins = append(plugins, actionPlugin)
	}

	return plugins, nil
}
