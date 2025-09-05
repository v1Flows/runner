package platform

import (
	"github.com/v1Flows/runner/config"
)

func GetPlatformConfig(cfg *config.Config) (string, string, string) {
	configManager := config.GetInstance()

	if cfg == nil {
		cfg = configManager.GetConfig()
	}

	return cfg.ExFlow.URL, cfg.ExFlow.APIKey, cfg.ExFlow.RunnerID
}

func GetPlatformConfigPlain(cfg *config.Config) (string, string) {
	configManager := config.GetInstance()

	if cfg == nil {
		cfg = configManager.GetConfig()
	}

	return cfg.ExFlow.URL, cfg.ExFlow.APIKey
}
