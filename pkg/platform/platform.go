package platform

import (
	"strings"

	"github.com/v1Flows/runner/config"
)

func GetPlatformConfig(platform string, cfg config.Config) (string, string, string) {
	configManager := config.GetInstance()

	switch strings.ToLower(platform) {
	case "alertflow":

		var api_key string
		if cfg.Alertflow.APIKey == "" && cfg.Runner.SharedRunnerSecret != "" {
			api_key = cfg.Runner.SharedRunnerSecret
		} else {
			api_key = cfg.Alertflow.APIKey
		}

		return cfg.Alertflow.URL, api_key, configManager.GetRunnerID("alertflow")
	case "exflow":

		var api_key string
		if cfg.ExFlow.APIKey == "" && cfg.Runner.SharedRunnerSecret != "" {
			api_key = cfg.Runner.SharedRunnerSecret
		} else {
			api_key = cfg.ExFlow.APIKey
		}

		return cfg.ExFlow.URL, api_key, configManager.GetRunnerID("exflow")
	default:
		return "unknown_platform", "unknown_platform", "unknown_platform"
	}
}

func GetPlatformConfigPlain(platform string, cfg config.Config) (string, string) {
	switch strings.ToLower(platform) {
	case "alertflow":
		return cfg.Alertflow.URL, cfg.Alertflow.APIKey
	case "exflow":
		return cfg.ExFlow.URL, cfg.ExFlow.APIKey
	default:
		return "unknown_platform", "unknown_platform"
	}
}
