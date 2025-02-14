package executions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetSteps(executionID string) ([]bmodels.ExecutionSteps, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	url := cfg.Alertflow.URL + "/api/v1/executions/" + executionID + "/steps"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return []bmodels.ExecutionSteps{}, err
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return []bmodels.ExecutionSteps{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get step data from API: %s", url)
		err = fmt.Errorf("failed to get step data from API: %s", url)
		return []bmodels.ExecutionSteps{}, err
	}

	log.Debugf("Step data received from API: %s", url)

	var steps models.IncomingExecutionSteps
	err = json.NewDecoder(resp.Body).Decode(&steps)
	if err != nil {
		log.Fatal(err)
		return []bmodels.ExecutionSteps{}, err
	}

	return steps.StepsData, nil
}
