package platform

import (
	"strings"

	"github.com/v1Flows/runner/config"
)

func GetPlatformConfig(platform string, cfg *config.Config) (string, string, string) {
	configManager := config.GetInstance()

	if cfg == nil {
		cfg = configManager.GetConfig()
	}

	switch strings.ToLower(platform) {
	case "alertflow":
		return cfg.Alertflow.URL, cfg.Alertflow.APIKey, cfg.Alertflow.RunnerID
	case "exflow":
		return cfg.ExFlow.URL, cfg.ExFlow.APIKey, cfg.ExFlow.RunnerID
	default:
		return "unknown_platform", "unknown_platform", "unknown_platform"
	}
}

func GetPlatformConfigPlain(platform string, cfg *config.Config) (string, string) {
	configManager := config.GetInstance()

	if cfg == nil {
		cfg = configManager.GetConfig()
	}

	switch strings.ToLower(platform) {
	case "alertflow":
		return cfg.Alertflow.URL, cfg.Alertflow.APIKey
	case "exflow":
		return cfg.ExFlow.URL, cfg.ExFlow.APIKey
	default:
		return "unknown_platform", "unknown_platform"
	}
}
