package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/internal/common"
	"github.com/AlertFlow/runner/internal/endpoints"
	"github.com/AlertFlow/runner/internal/runner"
	"github.com/AlertFlow/runner/pkg/plugins"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"

	"github.com/alecthomas/kingpin/v2"
)

var (
	log        = logrus.New()
	version    = "1.0.0-beta8"
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

	log.Info("Starting AlertFlow Runner. Version: ", version)

	log.Info("Loading config")
	configManager := config.GetInstance()
	err := configManager.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg := configManager.GetConfig()

	logging(cfg.LogLevel)

	loadedPlugins, modelPlugins, actionPlugins, endpointPlugins := plugins.Init(cfg)

	actions := common.RegisterActions(actionPlugins)
	endpoints := endpoints.RegisterEndpoints(endpointPlugins)

	runner.RegisterAtAPI(version, modelPlugins, actions, endpoints)

	// RunnerID might have changed after registration, so fetch the config again
	cfg = configManager.GetConfig()

	go runner.SendHeartbeat()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	Init(cfg, router, actions, endpointPlugins, loadedPlugins)

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info("Shutting down...")
	plugins.ShutdownPlugins()
	log.Info("Shutdown complete")
}

func Init(cfg config.Config, router *gin.Engine, actions []models.Actions, endpointPlugins []models.Plugins, loadedPlugins map[string]plugins.Plugin) {
	switch strings.ToLower(cfg.Mode) {
	case "master":
		log.Info("Runner is in Master Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker(cfg, actions, loadedPlugins)
		log.Info("Starting Alert Listener")
		go endpoints.InitAlertRouter(cfg, router, endpointPlugins, loadedPlugins)
		go endpoints.ReadyEndpoint(cfg, router)
	case "worker":
		log.Info("Runner is in Worker Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker(cfg, actions, loadedPlugins)
		go endpoints.ReadyEndpoint(cfg, router)
	case "listener":
		log.Info("Runner is in Listener Mode")
		log.Info("Starting Alert Listener")
		go endpoints.InitAlertRouter(cfg, router, endpointPlugins, loadedPlugins)
		go endpoints.ReadyEndpoint(cfg, router)
	}
}
