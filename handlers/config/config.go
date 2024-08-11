package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
)

var Config RestfulConf

type AlertflowConf struct {
	URL    string `json:"url"`
	APIKey string `json:"apikey"`
}

type PayloadsConf struct {
	Enabled  bool     `json:",default=true"`
	Port     int      `json:"port"`
	Managers []string `json:"managers"`
}

type PluginConf struct {
	Enable bool     `json:",default=false"`
	List   []string `json:"list"`
}

type RestfulConf struct {
	LogLevel  string `json:",default=Info"`
	RunnerID  string
	Alertflow AlertflowConf
	Payloads  PayloadsConf
	Plugins   PluginConf
}

func ReadConfig(configFile string) (*RestfulConf, error) {
	if err := conf.Load(configFile, &Config); err != nil {
		// check if we have os env vars
		if os.Getenv("ALERTFLOW_URL") != "" && os.Getenv("ALERTFLOW_APIKEY") != "" && os.Getenv("RUNNER_ID") != "" {
			Config.Alertflow.URL = os.Getenv("ALERTFLOW_URL")
			Config.Alertflow.APIKey = os.Getenv("ALERTFLOW_APIKEY")
			Config.RunnerID = os.Getenv("RUNNER_ID")

			if os.Getenv("PLUGIN_ENABLE") != "" {
				Config.Plugins.Enable, _ = strconv.ParseBool(os.Getenv("PLUGIN_ENABLE"))
			}

			return &Config, nil
		}
		log.Fatal("Error Loading Config File: ", err)
	}
	log.Info("Loaded Config File: ", configFile)

	return &Config, nil
}
