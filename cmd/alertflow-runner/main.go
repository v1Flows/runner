package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/internal/common"
	payloadendpoints "github.com/AlertFlow/runner/internal/payload_endpoints"
	"github.com/AlertFlow/runner/internal/runner"
	"github.com/AlertFlow/runner/pkg/plugins"

	"github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
)

const version string = "0.21.0"

var (
	configFile = kingpin.Flag("config", "Config File").Short('c').Default("config.yaml").String()
)

func logging(logLevel string) {
	logLevel = strings.ToLower(logLevel)

	if logLevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "warn" {
		log.SetLevel(log.WarnLevel)
	} else if logLevel == "error" {
		log.SetLevel(log.ErrorLevel)
	} else if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
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

	// result, err := loadedPlugins["alertmanager"].Execute(map[string]string{"target": "example.com"})
	// if err != nil {
	// 	log.Fatalf("Error executing plugin %s: %v", "test", err)
	// }

	// fmt.Printf("Plugin %s Execute Result: %s\n", "test", result)

	actions := common.RegisterActions(actionPlugins)
	endpoints := payloadendpoints.RegisterEndpoints(endpointPlugins)

	go payloadendpoints.InitPayloadRouter(cfg.PayloadEndpoints.Port, endpointPlugins, loadedPlugins)

	runner.RegisterAtAPI(version, modelPlugins, actions, endpoints)
	go runner.SendHeartbeat()

	Init(cfg)

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info("Shutting down...")
	plugins.ShutdownPlugins()
	log.Info("Shutdown complete")
}

func Init(cfg config.Config) {
	switch strings.ToLower(cfg.Mode) {
	case "master":
		log.Info("Runner is in Master Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker()
		log.Info("Starting Payload Listener")
		// go payloadhandler.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	case "worker":
		log.Info("Runner is in Worker Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker()
	case "listener":
		log.Info("Runner is in Listener Mode")
		log.Info("Starting Payload Listener")
		// go payloadhandler.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	}
}
