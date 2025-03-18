package platform

import (
	"strings"

	"github.com/v1Flows/runner/config"
)

func GetPlatformConfig(platform string, cfg config.Config) (string, string, string) {
	configManager := config.GetInstance()

	switch strings.ToLower(platform) {
	case "alertflow":
		return cfg.Alertflow.URL, cfg.Alertflow.APIKey, configManager.GetRunnerID("alertflow")
	case "exflow":
		return cfg.ExFlow.URL, cfg.ExFlow.APIKey, configManager.GetRunnerID("exflow")
	default:
		return "", "", ""
	}
}

func GetPlatformConfigPlain(platform string, cfg config.Config) (string, string) {
	switch strings.ToLower(platform) {
	case "alertflow":
		return cfg.Alertflow.URL, cfg.Alertflow.APIKey
	case "exflow":
		return cfg.ExFlow.URL, cfg.ExFlow.APIKey
	default:
		return "what", "what"
	}
}
