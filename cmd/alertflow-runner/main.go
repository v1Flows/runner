package main

import (
	"strings"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/internal/common"
	payloadendpoints "github.com/AlertFlow/runner/internal/payload_endpoints"
	"github.com/AlertFlow/runner/internal/runner"
	"github.com/AlertFlow/runner/pkg/plugin"

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

	manager, plugins, actions, payloadEndpoints := plugin.Init(cfg)

	common.RegisterActions(actions)
	go payloadendpoints.InitPayloadRouter(cfg.PayloadEndpoints.Port, manager, plugins, payloadEndpoints)

	runner.RegisterAtAPI(version, plugins, actions, payloadEndpoints)
	go runner.SendHeartbeat()

	Init(manager, cfg)

	<-make(chan struct{})

	defer manager.Cleanup()
}

func Init(manager *plugin.Manager, cfg config.Config) {
	switch strings.ToLower(cfg.Mode) {
	case "master":
		log.Info("Runner is in Master Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker(manager)
		log.Info("Starting Payload Listener")
		// go payloadhandler.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	case "worker":
		log.Info("Runner is in Worker Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker(manager)
	case "listener":
		log.Info("Runner is in Listener Mode")
		log.Info("Starting Payload Listener")
		// go payloadhandler.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	}
}
