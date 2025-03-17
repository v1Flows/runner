package executions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/models"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetSteps(cfg config.Config, executionID string) ([]shared_models.ExecutionSteps, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	platform, ok := platform.GetPlatformForExecution(executionID)
	if !ok {
		log.Error("Failed to get platform")
		return []shared_models.ExecutionSteps{}, fmt.Errorf("failed to get platform")
	}

	url, apiKey, _ := common.GetPlatformConfig(platform, cfg)

	parsedUrl := url + "/api/v1/executions/" + executionID + "/steps"
	req, err := http.NewRequest("GET", parsedUrl, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return []shared_models.ExecutionSteps{}, err
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return []shared_models.ExecutionSteps{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get step data from %s API: %s", platform, url)
		err = fmt.Errorf("failed to get step data from %s API: %s", platform, url)
		return []shared_models.ExecutionSteps{}, err
	}

	log.Debugf("Step data received from %s API: %s", platform, url)

	var steps models.IncomingExecutionSteps
	err = json.NewDecoder(resp.Body).Decode(&steps)
	if err != nil {
		log.Fatal(err)
		return []shared_models.ExecutionSteps{}, err
	}

	return steps.StepsData, nil
}
