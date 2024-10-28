package executions

import (
	"alertflow-runner/config"
	"alertflow-runner/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetSteps(executionID string) ([]models.ExecutionSteps, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := config.Config.Alertflow.URL + "/api/v1/executions/" + executionID + "/steps"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return []models.ExecutionSteps{}, err
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return []models.ExecutionSteps{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get step data from API: %s", url)
		err = fmt.Errorf("failed to get step data from API: %s", url)
		return []models.ExecutionSteps{}, err
	}

	log.Debugf("Step data received from API: %s", url)

	var steps models.IncomingExecutionSteps
	err = json.NewDecoder(resp.Body).Decode(&steps)
	if err != nil {
		log.Fatal(err)
		return []models.ExecutionSteps{}, err
	}

	return steps.StepsData, nil
}
