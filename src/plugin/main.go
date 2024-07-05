package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
)

func InitializePlugins(pluginList []string) {
	const pluginDir = "plugins"

	// Create the plugin directory
	if err := os.Mkdir(pluginDir, 0755); err != nil && !os.IsExist(err) {
		log.Fatal("Failed to create plugin directory: ", err)
	}

	// Iterate over the list of plugins
	for _, plugin := range pluginList {
		pluginDir := filepath.Join(pluginDir, plugin)

		// Clone the plugin
		if _, err := git.PlainClone(pluginDir, false, &git.CloneOptions{
			URL: "https://" + plugin,
		}); err != nil {
			if err != git.ErrRepositoryAlreadyExists {
				log.Fatalf("Failed to clone plugin %s: %v", plugin, err)
			}
		}

		// Initialize the plugin
		if err := initPlugin(pluginDir); err != nil {
			log.Fatalf("Failed to initialize plugin %s: %v", plugin, err)
		}
	}
}

func initPlugin(pluginDir string) error {
	// Change into the plugin directory
	if err := os.Chdir(pluginDir); err != nil {
		return fmt.Errorf("failed to change into plugin directory: %v", err)
	}

	// Vendor go modules
	if err := exec.Command("go", "mod", "vendor").Run(); err != nil {
		return fmt.Errorf("failed to vendor go modules: %v", err)
	}

	// Get go packages
	if err := exec.Command("go", "get").Run(); err != nil {
		return fmt.Errorf("failed to get go packages: %v", err)
	}

	// Run the plugin
	if err := exec.Command("go", "run", "main.go").Run(); err != nil {
		return fmt.Errorf("failed to run plugin: %v", err)
	}

	return nil
}
