package common

import "github.com/v1Flows/runner/config"

func GetPlatformConfig(platform string, cfg config.Config) (string, string, string) {
	switch platform {
	case "alertflow":
		return cfg.Alertflow.URL, cfg.Alertflow.APIKey, cfg.Alertflow.RunnerID
	case "exflow":
		return cfg.ExFlow.URL, cfg.ExFlow.APIKey, cfg.ExFlow.RunnerID
	default:
		return "", "", ""
	}
}
