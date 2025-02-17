package common

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/plugins"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

type IncomingExecutions struct {
	Executions []bmodels.Executions `json:"executions"`
}

var executionsExecuted = make(map[string]bool)

func StartWorker(cfg config.Config, actions []models.Actions, loadedPlugins map[string]plugins.Plugin) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := cfg.Alertflow.URL + "/api/v1/runners/" + cfg.Alertflow.RunnerID + "/executions/pending"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		var resp *http.Response
		var err error
		for i := 0; i < 3; i++ {
			resp, err = client.Do(req)
			if err != nil {
				log.Errorf("Failed to send request: %v", err)
				time.Sleep(5 * time.Second) // Add delay before retrying
				continue
			}

			if resp.StatusCode != 200 {
				log.Errorf("Failed to get waiting executions from API: %s, attempt %d", url, i+1)
				time.Sleep(5 * time.Second) // Add delay before retrying
				continue
			}

			log.Debugf("Executions received from API: %s", url)

			var executions IncomingExecutions
			err = json.NewDecoder(resp.Body).Decode(&executions)
			resp.Body.Close() // Close the body after reading
			if err != nil {
				log.Errorf("Failed to decode response body: %v", err)
				time.Sleep(5 * time.Second) // Add delay before retrying
				continue
			}

			for _, execution := range executions.Executions {
				// Check if execution is already executed
				if _, ok := executionsExecuted[execution.ID.String()]; ok {
					continue
				}
				// Add execution to executed map
				executionsExecuted[execution.ID.String()] = true

				// Process one execution at a time
				startProcessing(cfg, actions, loadedPlugins, execution)
			}
			break
		}
		if resp.StatusCode != 200 {
			log.Fatalf("Failed to get waiting executions from API after 3 attempts: %s", url)
		}
	}
}
