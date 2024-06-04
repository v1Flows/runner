package main

import (
	"alertflow-runner/src/config"
	"alertflow-runner/src/heartbeat"
	"alertflow-runner/src/incoming"
	"alertflow-runner/src/plugin"
	"alertflow-runner/src/register"

	"github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
)

const version string = "1.0.0"

var (
	configFile = kingpin.Flag("config.file", "Config File").String()

	logLevel = kingpin.Flag("log.level", "Log Level").Default("Info").String()

	runnerID        = kingpin.Flag("runner.id", "Runner ID").String()
	alertflowURL    = kingpin.Flag("alertflow.url", "Alertflow URL").String()
	alertflowAPIKey = kingpin.Flag("alertflow.apikey", "Alertflow API Key").String()

	pluginEnable = kingpin.Flag("plugin.enable", "Plugin Enable").Bool()
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

var (
	ApiURL   = ""
	ApiKey   = ""
	RunnerID = ""
)

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

	ApiURL = config.Alertflow.URL
	ApiKey = config.Alertflow.APIKey
	RunnerID = config.RunnerID

	go register.RegisterAtAPI(config.Alertflow.URL, config.Alertflow.APIKey, config.RunnerID, version)
	go heartbeat.SendHeartbeat(config.Alertflow.URL, config.Alertflow.APIKey, config.RunnerID)

	if config.ReceivePayloads.Enabled {
		log.Info("Starting ReceivePayloads")
		go incoming.InitPayloadRouter(config.ReceivePayloads.Port)
	}

	if config.Plugins.Enable {
		log.Info("Starting Plugin")
		go plugin.InitRPCServer()
		go plugin.InitializePlugins(config.Plugins.List)
	}

	<-make(chan struct{})
}
