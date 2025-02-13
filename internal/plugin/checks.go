package plugin

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func isPluginPresent(pluginName string, pluginDir string) bool {
	_, err := os.Stat(pluginDir + "/" + pluginName + ".so")
	return !os.IsNotExist(err)
}

func isPluginUpToDate(pluginName string, pluginVersion string, pluginDir string) bool {
	versionsFile := filepath.Join(pluginDir, ".versions")

	file, err := os.Open(versionsFile)
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
