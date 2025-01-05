package runner

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func RegisterAtAPI(version string, plugins []models.Plugin, actions []models.ActionDetails, payloadInjectors []models.PayloadEndpoint) {
	register := models.Register{
		ID:            config.Config.Alertflow.RunnerID,
		Registered:    true,
		LastHeartbeat: time.Now(),
		Version:       version,
		Mode:          config.Config.Mode,
	}

	// Convert plugins to JSON
	pluginsJSON, err := json.Marshal(plugins)
	if err != nil {
		log.Fatal(err)
	}
	register.Plugins = json.RawMessage(pluginsJSON)

	// Convert actions to JSON
	actionsJSON, err := json.Marshal(actions)
	if err != nil {
		log.Fatal(err)
	}
	register.Actions = json.RawMessage(actionsJSON)

	// Convert payloadInjectors to JSON
	payloadInjectorsJSON, err := json.Marshal(payloadInjectors)
	if err != nil {
		log.Fatal(err)
	}
	register.PayloadEndpoints = json.RawMessage(payloadInjectorsJSON)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(register)
	req, err := http.NewRequest("PUT", config.Config.Alertflow.URL+"/api/v1/runners/register", payloadBuf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 201 {
		log.Error("Failed to register at AlertFlow")
		log.Error("Response: ", string(body))
		panic("Failed to register at AlertFlow")
	}

	var response struct {
		RunnerID string `json:"runner_id"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	config.UpdateRunnerID(response.RunnerID)

	log.Info("Runner registered at AlertFlow. ID: ", config.GetRunnerID())
}
