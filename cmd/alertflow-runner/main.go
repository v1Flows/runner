package main

import (
	"strings"

	"gitlab.justlab.xyz/alertflow-public/runner/config"
	"gitlab.justlab.xyz/alertflow-public/runner/internal/actions"
	"gitlab.justlab.xyz/alertflow-public/runner/internal/common"
	"gitlab.justlab.xyz/alertflow-public/runner/internal/runner"
	payloadhandler "gitlab.justlab.xyz/alertflow-public/runner/pkg/handlers/payload"

	"github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
)

const version string = "0.12.5-beta"

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
	config, err := config.ReadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	logging(config.LogLevel)

	actions := actions.Init()
	payloadInjectors := payloadhandler.Init()

	go runner.RegisterAtAPI(version, actions, payloadInjectors)
	go runner.SendHeartbeat()

	Init()

	<-make(chan struct{})
}

func Init() {
	switch strings.ToLower(config.Config.Mode) {
	case "master":
		log.Info("Runner is in Master Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker(config.Config.Alertflow.URL, config.Config.Alertflow.APIKey, config.Config.RunnerID)
		log.Info("Starting Payload Listener")
		go payloadhandler.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	case "worker":
		log.Info("Runner is in Worker Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker(config.Config.Alertflow.URL, config.Config.Alertflow.APIKey, config.Config.RunnerID)
	case "listener":
		log.Info("Runner is in Listener Mode")
		log.Info("Starting Payload Listener")
		go payloadhandler.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	}
}
