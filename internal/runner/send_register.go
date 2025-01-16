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

	for i := 0; i < 3; i++ {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Errorf("Failed to send request: %v", err)
			time.Sleep(5 * time.Second) // Add delay before retrying
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close the body after reading
		if err != nil {
			log.Errorf("Failed to read response body: %v", err)
			time.Sleep(5 * time.Second) // Add delay before retrying
			continue
		}

		if resp.StatusCode == 201 {
			var response struct {
				RunnerID string `json:"runner_id"`
			}
			if err := json.Unmarshal(body, &response); err != nil {
				log.Fatal(err)
			}

			runner_id := ""
			if response.RunnerID == "" {
				runner_id = config.Config.Alertflow.RunnerID
			} else {
				runner_id = response.RunnerID
			}

			config.UpdateRunnerID(runner_id)

			log.Info("Runner registered at AlertFlow. ID: ", config.GetRunnerID())
			return
		} else {
			log.Errorf("Failed to register at AlertFlow, attempt %d", i+1)
			log.Errorf("Response: %s", string(body))
			time.Sleep(5 * time.Second) // Add delay before retrying
		}
	}
	log.Fatalf("Failed to register at AlertFlow after 3 attempts")
}
