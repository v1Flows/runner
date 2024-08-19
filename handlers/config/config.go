package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
)

var Config RestfulConf

type AlertflowConf struct {
	URL    string `json:",env=ALERTFLOW_URL"`
	APIKey string `json:",env=ALERTFLOW_API_KEY"`
}

type PayloadsConf struct {
	Enabled  bool     `json:",env=PAYLOADS_ENABLED,default=true"`
	Port     int      `json:",env=PAYLOADS_PORT,default=8080"`
	Managers []string `json:",env=PAYLOADS_MANAGERS,default=['alertmanager']"`
}

type RestfulConf struct {
	LogLevel  string `json:",env=LOG_LEVEL"`
	RunnerID  string `json:",env=RUNNER_ID"`
	Mode      string `json:",env=MODE"`
	Alertflow AlertflowConf
	Payloads  PayloadsConf
}

func ReadConfig(configFile string) (*RestfulConf, error) {
	conf.MustLoad(configFile, &Config, conf.UseEnv())
	log.Info("Loaded Config File: ", configFile)

	return &Config, nil
}
