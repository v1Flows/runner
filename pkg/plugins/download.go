// filepath: /Users/Justin.Neubert/projects/v1flows/alertflow/runner/pkg/plugins/download.go
package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlertFlow/runner/config"
	log "github.com/sirupsen/logrus"
)

// DownloadAndBuildPlugins downloads and builds plugins from GitHub
func DownloadAndBuildPlugins(pluginRepos []config.PluginConfig, buildDir string, pluginDir string) (map[string]string, error) {
	pluginPaths := make(map[string]string)

	// Create the build directory if it doesn't exist
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		err := os.MkdirAll(buildDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create build directory: %v", err)
		}
	}

	// Create the pluginDir directory if it doesn't exist
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		err := os.MkdirAll(pluginDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create build directory: %v", err)
		}
	}

	for _, plugin := range pluginRepos {
		// Define the plugin path
		pluginPath := filepath.Join(pluginDir, plugin.Name)

		// Clone the plugin repository
		log.Info("Cloning plugin ", plugin.Name)
		repoDir := filepath.Join(buildDir, plugin.Name)
		cmd := exec.Command("git", "clone", plugin.Repository, repoDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to clone plugin %s: %v\nOutput: %s", plugin.Name, err, string(output))
		}

		// Check out the specified version if provided
		if plugin.Version != "" {
			cmd = exec.Command("git", "checkout", plugin.Version)
			cmd.Dir = repoDir
			output, err = cmd.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("failed to checkout version %s for plugin %s: %v\nOutput: %s", plugin.Version, plugin.Name, err, string(output))
			}
		}

		// Build the plugin
		log.Info("Building plugin ", plugin.Name)
		cmd = exec.Command("go", "build", "-o", pluginPath)
		cmd.Env = append(os.Environ(), "GO111MODULE=on")
		cmd.Dir = repoDir
		output, err = cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to build plugin %s: %v\nOutput: %s", plugin.Name, err, string(output))
		}

		pluginPaths[plugin.Name] = pluginPath
	}

	// remove the buildDir directory
	err := os.RemoveAll(buildDir)
	if err != nil {
		return nil, fmt.Errorf("failed to remove build directory: %v", err)
	}

	return pluginPaths, nil
}
