package plugins

import (
	"fmt"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func CloneAndBuildPlugin(repoURL, pluginDir string, pluginRawRepos string, pluginName string) error {
	// Clone the repository
	cmd := exec.Command("git", "clone", "https://"+repoURL, pluginRawRepos)
	if err := cmd.Run(); err != nil {
		log.Error("failed to clone repository: " + err.Error())
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Build the plugin
	cmd = exec.Command("go", "build", "-buildmode=plugin", "-o", filepath.Join(pluginDir, pluginName+".so"), pluginRawRepos+"/plugin.go")
	if err := cmd.Run(); err != nil {
		log.Error("failed to build plugin: " + err.Error())
		return fmt.Errorf("failed to build plugin: %w", err)
	}

	return nil
}
