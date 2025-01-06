package config

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

const (
	defaultLogLevel = "info"
	defaultMode     = "master"
	defaultPort     = 8081
	defaultRunnerID = ""
	defaultAPIKey   = ""
)

var (
	Config RestfulConf
	mu     sync.RWMutex
)

type AlertflowConf struct {
	URL      string `json:"url"`
	RunnerID string `json:"runnerID,omitempty"`
	APIKey   string `json:"apiKey,omitempty"`
}

type PayloadEndpointsConf struct {
	Port int `json:"port,omitempty"`
}

type PluginConf struct {
	Name    string `json:"name,omitempty"`
	Url     string `json:"url,omitempty"`
	Version string `json:"version,omitempty"`
}

type RestfulConf struct {
	LogLevel         string               `json:"logLevel,omitempty"`
	Mode             string               `json:"mode,omitempty"`
	Alertflow        AlertflowConf        `json:"alertflow"`
	PayloadEndpoints PayloadEndpointsConf `json:"payload_endpoints"`
	Plugins          []PluginConf         `json:"plugins"`
}

func (c *RestfulConf) SetDefaults() {
	if c.LogLevel == "" {
		c.LogLevel = defaultLogLevel
	}
	if c.Mode == "" {
		c.Mode = defaultMode
	}
	if c.PayloadEndpoints.Port == 0 {
		c.PayloadEndpoints.Port = defaultPort
	}
}

func (c *RestfulConf) Validate() error {
	if c.LogLevel == "" {
		return fmt.Errorf("log level is required")
	}
	if c.Mode == "" {
		return fmt.Errorf("mode is required")
	}
	if c.PayloadEndpoints.Port == 0 {
		return fmt.Errorf("payload endpoints port is required")
	}
	if c.Alertflow.APIKey == "" {
		return fmt.Errorf("api key is required")
	}
	return nil
}

func ReadConfig(configFile string) (*RestfulConf, error) {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	log.Infoln("Loaded Config File:", configFile)

	if err := viper.Unmarshal(&Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	Config.SetDefaults()

	if err := Config.Validate(); err != nil {
		return nil, err
	}

	return &Config, nil
}

func UpdateRunnerID(runnerID string) {
	mu.Lock()
	defer mu.Unlock()
	Config.Alertflow.RunnerID = runnerID
}

func GetRunnerID() string {
	mu.RLock()
	defer mu.RUnlock()
	return Config.Alertflow.RunnerID
}
