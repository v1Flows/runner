package plugin

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func removePluginFromVersionsFile(pluginName string, pluginDir string) error {
	versionsFile := filepath.Join(pluginDir, ".versions")

	// remove plugin from versions file
	file, err := os.Open(versionsFile)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to open .versions file: %v", err)
	}
	defer file.Close()

	lines := []string{}
	found := false
	for {
		var name, version string
		_, err := fmt.Fscanf(file, "%s %s\n", &name, &version)
		if err != nil {
			break
		}
		if name == pluginName {
			found = true
			continue
		}
		lines = append(lines, fmt.Sprintf("%s %s", name, version))
	}

	if !found {
		return nil
	}

	file, err = os.Create(versionsFile)
	if err != nil {
		return fmt.Errorf("failed to create .versions file: %v", err)
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to .versions file: %v", err)
		}
	}

	return nil
}

func updatePluginVersion(pluginName string, pluginVersion string, pluginDir string) error {
	versionsFile := filepath.Join(pluginDir, ".versions")
	// Read the existing .versions file
	file, err := os.Open(versionsFile)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(versionsFile)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		if parts[0] != pluginName {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write the updated content back to the .versions file
	file, err = os.OpenFile(versionsFile, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// Append the new version entry
	_, err = file.WriteString(fmt.Sprintf("%s %s\n", pluginName, pluginVersion))
	return err
}
