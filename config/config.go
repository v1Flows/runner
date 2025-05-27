package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ConfigurationManager handles all configuration operations
type ConfigurationManager struct {
	config *Config
	mu     sync.RWMutex
	viper  *viper.Viper
}

// Config represents the application configuration
type Config struct {
	LogLevel     string            `mapstructure:"log_level" validate:"required,oneof=debug info warn error"`
	Mode         string            `mapstructure:"mode" validate:"required,oneof=master worker"`
	Alertflow    AlertflowConfig   `mapstructure:"alertflow" validate:"required"`
	ExFlow       exflowConfig      `mapstructure:"exflow" validate:"required"`
	ApiEndpoint  ApiEndpointConfig `mapstructure:"api_endpoint" validate:"required"`
	WorkspaceDir string            `mapstructure:"workspace_dir" validate:"dir"`
	PluginDir    string            `mapstructure:"plugin_dir" validate:"dir"`
	Plugins      []PluginConfig    `mapstructure:"plugins"`
	Runner       RunnerConf        `mapstructure:"runner"`
}

type AlertflowConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	URL      string `mapstructure:"url" validate:"required,url"`
	RunnerID string `mapstructure:"runner_id"`
	APIKey   string `mapstructure:"api_key"`
}

type exflowConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	URL      string `mapstructure:"url" validate:"required,url"`
	RunnerID string `mapstructure:"runner_id"`
	APIKey   string `mapstructure:"api_key"`
}

type ApiEndpointConfig struct {
	Port int `mapstructure:"port" validate:"required,min=1024,max=65535"`
}

type PluginConfig struct {
	Name    string `mapstructure:"name" validate:"required"`
	Url     string `mapstructure:"url" validate:"required,url"`
	Version string `mapstructure:"version" validate:"required"`
}

type RunnerConf struct {
	SharedRunnerSecret string `mapstructure:"shared_runner_secret"`
}

const (
	defaultLogLevel = "info"
	defaultMode     = "master"
	defaultPort     = 8081
)

var (
	instance *ConfigurationManager
	once     sync.Once
)

// GetInstance returns the singleton configuration manager instance
func GetInstance() *ConfigurationManager {
	once.Do(func() {
		instance = &ConfigurationManager{
			viper: viper.New(),
		}
	})
	return instance
}

// LoadConfig initializes the configuration from file and environment
func (cm *ConfigurationManager) LoadConfig(configFile string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Set up Viper
	cm.viper.SetConfigFile(configFile)
	cm.viper.SetEnvPrefix("RUNNER")
	cm.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cm.viper.AutomaticEnv()

	// Bind specific environment variables
	envBindings := map[string]string{
		"alertflow.api_key": "RUNNER_ALERTFLOW_API_KEY",
		"exflow.api_key":    "RUNNER_EXFLOW_API_KEY",
		"plugin_dir":        "RUNNER_PLUGIN_DIR",
	}

	for configKey, envVar := range envBindings {
		if err := cm.viper.BindEnv(configKey, envVar); err != nil {
			return fmt.Errorf("failed to bind env var %s: %w", envVar, err)
		}
	}

	// Read configuration file
	if err := cm.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Create new config instance
	var config Config

	// Set defaults
	cm.setDefaults(&config)

	// Unmarshal configuration
	if err := cm.viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := cm.validateConfig(&config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Store the config
	cm.config = &config

	log.WithFields(log.Fields{
		"file":    configFile,
		"content": cm.viper.AllSettings(),
	}).Debug("Configuration loaded successfully")

	return nil
}

func (cm *ConfigurationManager) setDefaults(config *Config) {
	if config.LogLevel == "" {
		config.LogLevel = defaultLogLevel
	}
	if config.Mode == "" {
		config.Mode = defaultMode
	}
	if config.ApiEndpoint.Port == 0 {
		config.ApiEndpoint.Port = defaultPort
	}
	if config.WorkspaceDir == "" {
		// get the current working directory and add plugins folder
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current working directory: %v", err)
		}
		config.WorkspaceDir = currentDir + "/workspace"
	}
	if config.PluginDir == "" {
		// get the current working directory and add plugins folder
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current working directory: %v", err)
		}
		config.PluginDir = currentDir + "/plugins"
	}
	config.Alertflow.Enabled = true
	config.ExFlow.Enabled = true
}

func (cm *ConfigurationManager) validateConfig(config *Config) error {
	if config.Alertflow.Enabled {
		if config.Alertflow.APIKey == "" && config.Runner.SharedRunnerSecret == "" {
			return fmt.Errorf("alertflow.api_key or runner.shared_runner_secret is required")
		}
		if config.Alertflow.URL == "" {
			return fmt.Errorf("alertflow URL is required")
		}
	}
	if config.ExFlow.Enabled {
		if config.ExFlow.APIKey == "" && config.Runner.SharedRunnerSecret == "" {
			return fmt.Errorf("exflow.api_key or runner.shared_runner_secret is required")
		}
		if config.ExFlow.URL == "" {
			return fmt.Errorf("exflow URL is required")
		}
	}

	return nil
}

// GetConfig returns a copy of the current configuration
func (cm *ConfigurationManager) GetConfig() Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return *cm.config
}

// UpdateRunnerID updates the runner ID in the configuration for both Alertflow and ExFlow
func (cm *ConfigurationManager) UpdateRunnerID(platform, id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if platform == "alertflow" {
		cm.config.Alertflow.RunnerID = id
	}
	if platform == "exflow" {
		cm.config.ExFlow.RunnerID = id
	}
}

// UpdateRunnerApiKey updates the runner api_key in the configuration for both Alertflow and ExFlow
func (cm *ConfigurationManager) UpdateRunnerApiKey(platform, apiKey string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if platform == "alertflow" {
		cm.config.Alertflow.APIKey = apiKey
	}
	if platform == "exflow" {
		cm.config.ExFlow.APIKey = apiKey
	}
}

// GetRunnerIDs returns the current runner IDs for both Alertflow and ExFlow
func (cm *ConfigurationManager) GetRunnerID(platform string) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if platform == "alertflow" {
		return cm.config.Alertflow.RunnerID
	}
	if platform == "exflow" {
		return cm.config.ExFlow.RunnerID
	}

	return ""
}

// GetRunnerIDs returns the current runner apiKey for both Alertflow and ExFlow
func (cm *ConfigurationManager) GetRunnerApiKey(platform string) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if platform == "alertflow" {
		return cm.config.Alertflow.APIKey
	}
	if platform == "exflow" {
		return cm.config.ExFlow.APIKey
	}

	return ""
}

// ReloadConfig reloads the configuration from the file
func (cm *ConfigurationManager) ReloadConfig() error {
	return cm.LoadConfig(cm.viper.ConfigFileUsed())
}
