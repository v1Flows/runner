package config

import (
	"os"
	"strconv"
	"strings"

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

type RestfulConf struct {
	LogLevel  string `json:",default=Info"`
	RunnerID  string
	Alertflow AlertflowConf
	Payloads  PayloadsConf
}

func ReadConfig(configFile string) (*RestfulConf, error) {
	if configFile == "" {
		Config.RunnerID = os.Getenv("RUNNER_ID")
		Config.Alertflow.URL = os.Getenv("ALERTFLOW_URL")
		Config.Alertflow.APIKey = os.Getenv("ALERTFLOW_APIKEY")

		Config.Payloads.Enabled, _ = strconv.ParseBool(os.Getenv("PAYLOADS_ENABLED"))
		Config.Payloads.Port, _ = strconv.Atoi(os.Getenv("PAYLOADS_PORT"))
		Config.Payloads.Managers = strings.Split(os.Getenv("PAYLOADS_MANAGERS"), ",")

		if Config.RunnerID == "" || Config.Alertflow.URL == "" || Config.Alertflow.APIKey == "" {
			log.Fatal("Missing Required Config Values")
		}

		return &Config, nil
	}
	if err := conf.Load(configFile, &Config); err != nil {
		log.Fatal("Error Loading Config File: ", err)
	}

	log.Info("Loaded Config File: ", configFile)

	return &Config, nil
}
