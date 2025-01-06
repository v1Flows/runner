package plugins

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func CloneAndBuildPlugin(repoURL, pluginDir string, pluginTempDir string, pluginName string, pluginVersion string) error {
	if repoURL == "" {
		repoURL = "https://github.com/Alertflow/rp-" + pluginName
	}

	err := prepareClone(pluginTempDir, pluginDir)
	if err != nil {
		log.Error("failed to prepare clone: ", err.Error())
		return err
	}

	outDir := filepath.Join(pluginTempDir, pluginName)

	// Clone the repository
	cmd := exec.Command("git", "clone", repoURL, "--branch", pluginVersion, outDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("failed to clone repository: %v, %s", err, stderr.String())
		return err
	}

	// Build the plugin
	cmd = exec.Command("go", "build", "-mod=mod", "-buildmode=plugin", "-o", filepath.Join(pluginDir, pluginName+".so"), outDir+"/plugin.go")
	if err := cmd.Run(); err != nil {
		log.Error("failed to build plugin: ", err.Error())
		return err
	}

	// Update the .versions file
	if err := UpdatePluginVersion(pluginName, pluginVersion); err != nil {
		log.Error("failed to update .versions file: ", err.Error())
		return err
	}

	return nil
}

func prepareClone(pluginTempDir string, pluginDir string) error {
	// Create the temporary directory
	if err := os.MkdirAll(pluginTempDir, 0755); err != nil {
		log.Error("failed to create temporary directory: " + err.Error())
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Create the plugin directory if it doesn't exist
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		log.Error("failed to create plugin directory: " + err.Error())
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	return nil
}
