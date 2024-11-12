package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gitlab.justlab.xyz/alertflow-public/runner/config"
	"gitlab.justlab.xyz/alertflow-public/runner/internal/common"
	"gitlab.justlab.xyz/alertflow-public/runner/internal/plugins"
	"gitlab.justlab.xyz/alertflow-public/runner/internal/runner"
	payloadhandler "gitlab.justlab.xyz/alertflow-public/runner/pkg/handlers/payload"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

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

func cloneAndBuildPlugin(repoURL, pluginDir string, pluginRawRepos string, pluginName string) error {
	// Clone the repository
	cmd := exec.Command("git", "clone", "https://"+repoURL, pluginRawRepos)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Build the plugin
	cmd = exec.Command("go", "build", "-buildmode=plugin", "-o", filepath.Join(pluginDir, pluginName+".so"), pluginRawRepos+"/"+pluginName+".go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build plugin: %w", err)
	}

	return nil
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

	pluginReposDir := "./temp/rawPlugins"
	pluginDir := "./plugins"
	os.MkdirAll(pluginReposDir, os.ModePerm)
	os.MkdirAll(pluginDir, os.ModePerm)

	for _, plugin := range config.Plugins {
		pluginPath := filepath.Join(pluginReposDir, plugin.Name)
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			log.Infof("Cloning and building plugin: %s", plugin.Name)
			if plugin.Version == "" {
				plugin.Version = "main"
			}
			err := plugins.CloneAndBuildPlugin(plugin.Url, pluginDir, pluginPath, plugin.Name, plugin.Version)
			if err != nil {
				log.Fatalf("Failed to clone and build plugin %s: %v", plugin.Name, err)
			}
		}
	}

	// cleanup the temp directory
	err = os.RemoveAll(pluginReposDir)
	if err != nil {
		log.Errorf("Failed to remove temp directory: %v", err)
	}

	plugins, err := plugins.LoadPlugins(pluginDir)
	if err != nil {
		log.Fatal(err)
	}

	actionsMap := make(map[string]models.ActionDetails)
	payloadEndpointsMap := make(map[string]models.PayloadEndpoint)
	for _, plugin := range plugins {
		p := plugin.Init()

		if p.Type == "action" {
			action := plugin.Details()
			actionsMap[action.Action.Type] = action.Action
			log.Infof("Loaded action plugin: %s", action.Action.Name)
		}
		if p.Type == "payload_endpoint" {
			payloadEndpoint := plugin.Details()
			payloadEndpointsMap[payloadEndpoint.Payload.Name] = payloadEndpoint.Payload
			log.Infof("Loaded payload endpoint plugin: %s", payloadEndpoint.Payload.Name)
		}
	}

	common.RegisterActions(actionsMap)

	payloadInjectors := payloadhandler.Init()

	actionsSlice := make([]models.ActionDetails, 0, len(actionsMap))
	for _, action := range actionsMap {
		actionsSlice = append(actionsSlice, action)
	}
	go runner.RegisterAtAPI(version, actionsSlice, payloadInjectors)
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
