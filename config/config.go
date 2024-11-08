package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
)

var Config RestfulConf

type AlertflowConf struct {
	URL    string `json:""`
	APIKey string `json:""`
}

type PayloadEndpointsConf struct {
	Port int `json:""`
}

type PluginConf struct {
	Name string `json:""`
	Url  string `json:""`
}

type RestfulConf struct {
	LogLevel         string `json:""`
	RunnerID         string `json:""`
	Mode             string `json:""`
	Alertflow        AlertflowConf
	PayloadEndpoints PayloadEndpointsConf
	Plugins          []PluginConf
}

func ReadConfig(configFile string) (*RestfulConf, error) {
	conf.MustLoad(configFile, &Config, conf.UseEnv())
	log.Info("Loaded Config File: ", configFile)

	return &Config, nil
}
