package runner

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"

	log "github.com/sirupsen/logrus"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
)

func RegisterAtAPI(platform string, version string, plugins []models.Plugins, actions []models.Actions, alertEndpoints []models.AlertEndpoints) {
	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	url, apiKey, runnerID := common.GetPlatformConfig(platform, cfg)

	var parsedRunnerID uuid.UUID
	var err error
	if runnerID != "" {
		parsedRunnerID, err = uuid.Parse(runnerID)
		if err != nil {
			log.Fatalf("Invalid RunnerID: %v", err)
		}
	}

	register := models.Runners{
		ID:             parsedRunnerID,
		Registered:     true,
		LastHeartbeat:  time.Now(),
		Version:        version,
		Mode:           cfg.Mode,
		Plugins:        plugins,
		Actions:        actions,
		AlertEndpoints: alertEndpoints,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(register)
	req, err := http.NewRequest("PUT", url+"/api/v1/runners/register", payloadBuf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", apiKey)

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
				runner_id = configManager.GetRunnerID(platform)
			} else {
				runner_id = response.RunnerID
			}

			configManager.UpdateRunnerID(platform, runner_id)

			log.Info("Runner registered at "+platform+". ID: ", configManager.GetRunnerID(platform))
			return
		} else {
			log.Errorf("Failed to register at "+platform+", attempt %d", i+1)
			log.Errorf("Response: %s", string(body))
			time.Sleep(5 * time.Second) // Add delay before retrying
		}
	}
	log.Fatalf("Failed to register at " + platform + " after 3 attempts")
}
