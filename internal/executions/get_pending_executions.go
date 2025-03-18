package internal_executions

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	af_models "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	ef_models "github.com/v1Flows/exFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/platform"
	platformfn "github.com/v1Flows/runner/pkg/platform"
	"github.com/v1Flows/runner/pkg/plugins"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

type IncomingSharedExecutions struct {
	Executions []shared_models.Executions `json:"executions"`
}

type IncomingAfExecutions struct {
	Executions []af_models.Executions `json:"executions"`
}

type IncomingEfExecutions struct {
	Executions []ef_models.Executions `json:"executions"`
}

func GetPendingExecutions(targetPlatform string, cfg config.Config, actions []shared_models.Action, loadedPlugins map[string]plugins.Plugin) {
	url, apiKey, runnerID := platform.GetPlatformConfig(targetPlatform, cfg)

	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	parsedUrl := url + "/api/v1/runners/" + runnerID + "/executions/pending"
	req, err := http.NewRequest("GET", parsedUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", apiKey)
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
				log.Errorf("Failed to get waiting executions from %s API: %s, attempt %d", targetPlatform, parsedUrl, i+1)
				time.Sleep(5 * time.Second) // Add delay before retrying
				continue
			}

			log.Debugf("Executions received from %s API: %s", targetPlatform, parsedUrl)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error(err)
				time.Sleep(5 * time.Second) // Add delay before retrying
				continue
			}
			resp.Body.Close() // Close the body after reading

			if targetPlatform == "alertflow" {
				var executions IncomingAfExecutions
				err := json.Unmarshal(body, &executions)
				if err != nil {
					log.Error(err)
					continue
				}

				var sharedExecutions IncomingSharedExecutions
				err = json.Unmarshal(body, &sharedExecutions)
				if err != nil {
					log.Error(err)
					continue
				}

				for index, execution := range executions.Executions {
					// Save platform information for the execution
					platformfn.SetPlatformForExecution(execution.ID.String(), targetPlatform)

					startProcessing(targetPlatform, cfg, actions, loadedPlugins, sharedExecutions.Executions[index], execution.AlertID)
				}
			}

			if targetPlatform == "exflow" {
				var executions IncomingSharedExecutions
				err := json.Unmarshal(body, &executions)
				if err != nil {
					log.Error(err)
					continue
				}

				for _, execution := range executions.Executions {
					// Save platform information for the execution
					platformfn.SetPlatformForExecution(execution.ID.String(), targetPlatform)

					startProcessing(targetPlatform, cfg, actions, loadedPlugins, execution, "")
				}
			}

		}
		if resp.StatusCode != 200 {
			log.Fatalf("Failed to get waiting executions from %s API after 3 attempts: %s", targetPlatform, parsedUrl)
		}
	}
}
