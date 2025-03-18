package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/endpoints"
	internal_executions "github.com/v1Flows/runner/internal/executions"
	"github.com/v1Flows/runner/internal/runner"
	"github.com/v1Flows/runner/internal/worker"
	"github.com/v1Flows/runner/pkg/plugins"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	"github.com/alecthomas/kingpin/v2"
)

var (
	log        = logrus.New()
	version    = "1.0.0-beta10"
	configFile = kingpin.Flag("config", "Path to configuration file").Short('c').String()
)

func logging(logLevel string) {
	logLevel = strings.ToLower(logLevel)

	if logLevel == "info" {
		log.SetLevel(logrus.InfoLevel)
	} else if logLevel == "warn" {
		log.SetLevel(logrus.WarnLevel)
	} else if logLevel == "error" {
		log.SetLevel(logrus.ErrorLevel)
	} else if logLevel == "debug" {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Info("Starting v1Flows Runner. Version: ", version)

	log.Info("Loading config")
	configManager := config.GetInstance()
	err := configManager.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg := configManager.GetConfig()

	logging(cfg.LogLevel)

	loadedPlugins, modelPlugins, actionPlugins, endpointPlugins := plugins.Init(cfg)

	actions := internal_executions.RegisterActions(actionPlugins)

	// RunnerID might have changed after registration, so fetch the config again
	cfg = configManager.GetConfig()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if cfg.Alertflow.Enabled {
		endpoints := endpoints.RegisterEndpoints(endpointPlugins)
		log.Info("Registering at AlertFlow")
		runner.RegisterAtAPI("alertflow", version, modelPlugins, actions, endpoints)
		go runner.SendHeartbeat("alertflow")
		Init("alertflow", cfg, router, actions, endpointPlugins, loadedPlugins)
	}

	if cfg.ExFlow.Enabled {
		log.Info("Registering at ExFlow")
		runner.RegisterAtAPI("exflow", version, modelPlugins, actions, nil)
		go runner.SendHeartbeat("exflow")
		Init("exflow", cfg, router, actions, endpointPlugins, loadedPlugins)
	}

	go endpoints.ReadyEndpoint(cfg, router)

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info("Shutting down...")
	plugins.ShutdownPlugins()
	log.Info("Shutdown complete")
}

func Init(platform string, cfg config.Config, router *gin.Engine, actions []shared_models.Action, endpointPlugins []shared_models.Plugin, loadedPlugins map[string]plugins.Plugin) {
	switch strings.ToLower(cfg.Mode) {
	case "master":
		log.Info("Runner is in Master Mode")
		log.Info("Starting Execution Checker")
		go worker.StartWorker(platform, cfg, actions, loadedPlugins)
		if platform == "alertflow" {
			log.Info("Starting Alert Listener")
			go endpoints.InitEndpointRouter(cfg, router, "alertflow", endpointPlugins, loadedPlugins)
		}
	case "worker":
		log.Info("Runner is in Worker Mode")
		log.Info("Starting Execution Checker")
		go worker.StartWorker(platform, cfg, actions, loadedPlugins)
	case "listener":
		log.Info("Runner is in Listener Mode")
		if platform == "alertflow" {
			log.Info("Starting Alert Listener")
			go endpoints.InitEndpointRouter(cfg, router, "alertflow", endpointPlugins, loadedPlugins)
		}
	}
}
