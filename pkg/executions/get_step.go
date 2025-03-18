package executions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/models"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetStep(cfg config.Config, executionID string, stepID string) (shared_models.ExecutionSteps, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	targetPlatform, ok := platform.GetPlatformForExecution(executionID)
	if !ok {
		log.Error("Failed to get platform")
		return shared_models.ExecutionSteps{}, fmt.Errorf("failed to get platform")
	}

	url, apiKey := platform.GetPlatformConfigPlain(targetPlatform, cfg)

	parsedUrl := url + "/api/v1/executions/" + executionID + "/steps/" + stepID
	req, err := http.NewRequest("GET", parsedUrl, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return shared_models.ExecutionSteps{}, err
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return shared_models.ExecutionSteps{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get step data from %s API: %s", targetPlatform, url)
		err = fmt.Errorf("failed to get step data from %s API: %s", targetPlatform, url)
		return shared_models.ExecutionSteps{}, err
	}

	log.Debugf("Step data received from %s API: %s", targetPlatform, url)

	var step models.IncomingExecutionStep
	err = json.NewDecoder(resp.Body).Decode(&step)
	if err != nil {
		log.Fatal(err)
		return shared_models.ExecutionSteps{}, err
	}

	return step.StepData, nil
}
