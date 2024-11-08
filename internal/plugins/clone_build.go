package plugins

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func CloneAndBuildPlugin(repoURL, pluginDir string, pluginRawRepos string, pluginName string, pluginVersion string) error {
	// Check if the plugin is already installed and up-to-date
	if isPluginUpToDate(pluginName, pluginVersion) {
		log.Infof("Plugin %s is already up-to-date", pluginName)
		return nil
	}

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

	// Update the .versions file
	if err := updatePluginVersion(pluginName, pluginVersion); err != nil {
		log.Error("failed to update .versions file: " + err.Error())
		return fmt.Errorf("failed to update .versions file: %w", err)
	}

	return nil
}

func isPluginUpToDate(pluginName, pluginVersion string) bool {
	file, err := os.Open(".versions")
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Error("failed to open .versions file: " + err.Error())
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		if parts[0] == pluginName && parts[1] == pluginVersion {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("failed to read .versions file: " + err.Error())
	}

	return false
}

func updatePluginVersion(pluginName, pluginVersion string) error {
	file, err := os.OpenFile(".versions", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s %s\n", pluginName, pluginVersion))
	return err
}
