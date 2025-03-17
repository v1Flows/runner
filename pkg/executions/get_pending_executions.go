package executions

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/alertflow"
	"github.com/v1Flows/runner/pkg/exflow"
	platformfn "github.com/v1Flows/runner/pkg/platform"
	"github.com/v1Flows/runner/pkg/plugins"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

type IncomingExecutions struct {
	Executions []shared_models.Executions `json:"executions"`
}

func GetPendingExecutions(platform string, cfg config.Config, actions []shared_models.Actions, loadedPlugins map[string]plugins.Plugin) {
	url, apiKey, runnerID := common.GetPlatformConfig(platform, cfg)

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
				log.Errorf("Failed to get waiting executions from %s API: %s, attempt %d", platform, parsedUrl, i+1)
				time.Sleep(5 * time.Second) // Add delay before retrying
				continue
			}

			log.Debugf("Executions received from %s API: %s", platform, parsedUrl)

			var executions IncomingExecutions
			err = json.NewDecoder(resp.Body).Decode(&executions)
			resp.Body.Close() // Close the body after reading
			if err != nil {
				log.Errorf("Failed to decode response body: %v", err)
				time.Sleep(5 * time.Second) // Add delay before retrying
				continue
			}

			for _, execution := range executions.Executions {
				// Save platform information for the execution
				platformfn.SetPlatformForExecution(execution.ID.String(), platform)

				if platform == "alertflow" {
					// Process one execution at a time
					alertflow.StartProcessing(platform, cfg, actions, loadedPlugins, execution)
				} else if platform == "exflow" {
					exflow.StartProcessing(platform, cfg, actions, loadedPlugins, execution)
				}
			}

		}
		if resp.StatusCode != 200 {
			log.Fatalf("Failed to get waiting executions from %s API after 3 attempts: %s", platform, parsedUrl)
		}
	}
}
