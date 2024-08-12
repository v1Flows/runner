package main

import (
	"alertflow-runner/handlers/actions"
	"alertflow-runner/handlers/config"
	"alertflow-runner/handlers/executions"
	"alertflow-runner/handlers/heartbeat"
	"alertflow-runner/handlers/payload"
	"alertflow-runner/handlers/register"

	"github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
)

const version string = "0.3.0-beta"

var (
	configFile = kingpin.Flag("config.file", "Config File").String()
)

func logging(logLevel string) {
	if logLevel == "Info" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "Warn" {
		log.SetLevel(log.WarnLevel)
	} else if logLevel == "Error" {
		log.SetLevel(log.ErrorLevel)
	} else if logLevel == "Debug" {
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

	go register.RegisterAtAPI(config.Alertflow.URL, config.Alertflow.APIKey, config.RunnerID, version, actions)
	go heartbeat.SendHeartbeat(config.Alertflow.URL, config.Alertflow.APIKey, config.RunnerID)

	Init()

	<-make(chan struct{})
}

func Init() {
	switch config.Config.Mode {
	case "master":
		log.Info("Runner is in Master Mode")
		log.Info("Starting Execution Checker")
		go executions.StartWorker(config.Config.Alertflow.URL, config.Config.Alertflow.APIKey, config.Config.RunnerID)
		log.Info("Starting Payload Listener")
		go payload.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	case "worker":
		log.Info("Runner is in Worker Mode")
		log.Info("Starting Execution Checker")
		go executions.StartWorker(config.Config.Alertflow.URL, config.Config.Alertflow.APIKey, config.Config.RunnerID)
	case "listener":
		log.Info("Runner is in Listener Mode")
		log.Info("Starting Payload Listener")
		go payload.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	}
}
