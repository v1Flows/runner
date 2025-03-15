package common

import "github.com/v1Flows/runner/config"

func GetPlatformConfig(platform string, cfg config.Config) (string, string, string) {
	configManager := config.GetInstance()

	switch platform {
	case "alertflow":
		return cfg.Alertflow.URL, cfg.Alertflow.APIKey, configManager.GetRunnerID("alertflow")
	case "exflow":
		return cfg.ExFlow.URL, cfg.ExFlow.APIKey, configManager.GetRunnerID("exflow")
	default:
		return "", "", ""
	}
}
