package platform

import (
	"strings"

	"github.com/v1Flows/runner/config"
)

func GetPlatformConfig(platform string, cfg config.Config) (string, string, string) {
	configManager := config.GetInstance()

	switch strings.ToLower(platform) {
	case "alertflow":
		return cfg.Alertflow.URL, configManager.GetRunnerApiKey("alertflow"), configManager.GetRunnerID("alertflow")
	case "exflow":
		return cfg.ExFlow.URL, configManager.GetRunnerApiKey("exflow"), configManager.GetRunnerID("exflow")
	default:
		return "unknown_platform", "unknown_platform", "unknown_platform"
	}
}

func GetPlatformConfigPlain(platform string, cfg config.Config) (string, string) {
	configManager := config.GetInstance()

	switch strings.ToLower(platform) {
	case "alertflow":
		return cfg.Alertflow.URL, configManager.GetRunnerApiKey("alertflow")
	case "exflow":
		return cfg.ExFlow.URL, configManager.GetRunnerApiKey("exflow")
	default:
		return "unknown_platform", "unknown_platform"
	}
}
