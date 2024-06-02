package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
)

type AlertflowConf struct {
	URL    string `json:"url"`
	APIKey string `json:"apikey"`
}

type PluginConf struct {
	Enable bool `json:",default=false"`
	Port   int  `json:",default=9854"`
}

type RestfulConf struct {
	LogLevel  string `json:",default=Info"`
	RunnerID  string
	Alertflow AlertflowConf
	Plugin    PluginConf
}

var config RestfulConf

func ReadConfig(configFile string) (*RestfulConf, error) {
	if err := conf.Load(configFile, &config); err != nil {
		// check if we have os env vars
		if os.Getenv("ALERTFLOW_URL") != "" && os.Getenv("ALERTFLOW_APIKEY") != "" && os.Getenv("RUNNER_ID") != "" {
			config.Alertflow.URL = os.Getenv("ALERTFLOW_URL")
			config.Alertflow.APIKey = os.Getenv("ALERTFLOW_APIKEY")
			config.RunnerID = os.Getenv("RUNNER_ID")

			if os.Getenv("PLUGIN_ENABLE") != "" {
				config.Plugin.Enable, _ = strconv.ParseBool(os.Getenv("PLUGIN_ENABLE"))
			}

			if os.Getenv("PLUGIN_PORT") != "" {
				port, _ := strconv.ParseInt(os.Getenv("PLUGIN_PORT"), 10, 32)
				config.Plugin.Port = int(port)
			}

			return &config, nil
		}
		log.Fatal("Error Loading Config File: ", err)
	}
	log.Info("Loaded Config File: ", configFile)

	return &config, nil
}

func GetConfig() RestfulConf {
	return config
}
