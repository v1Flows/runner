package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
)

var Config RestfulConf

type AlertflowConf struct {
	URL    string
	APIKey string
}

type PayloadsConf struct {
	Enabled  bool
	Port     int
	Managers []string
}

type RestfulConf struct {
	LogLevel  string
	RunnerID  string
	Mode      string
	Alertflow AlertflowConf
	Payloads  PayloadsConf
}

func ReadConfig(configFile string) (*RestfulConf, error) {
	if err := conf.Load(configFile, &Config); err != nil {
		log.Fatal("Error Loading Config File: ", err)
	}
	log.Info("Loaded Config File: ", configFile)

	return &Config, nil
}
