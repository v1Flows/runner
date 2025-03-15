package executions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetStep(cfg config.Config, executionID string, stepID string) (bmodels.ExecutionSteps, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	platform, ok := GetPlatformForExecution(executionID)
	if !ok {
		log.Error("Failed to get platform")
		return bmodels.ExecutionSteps{}, fmt.Errorf("failed to get platform")
	}

	url, apiKey, _ := common.GetPlatformConfig(platform, cfg)

	parsedUrl := url + "/api/v1/executions/" + executionID + "/steps/" + stepID
	req, err := http.NewRequest("GET", parsedUrl, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return bmodels.ExecutionSteps{}, err
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return bmodels.ExecutionSteps{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get step data from %s API: %s", platform, url)
		err = fmt.Errorf("failed to get step data from %s API: %s", platform, url)
		return bmodels.ExecutionSteps{}, err
	}

	log.Debugf("Step data received from %s API: %s", platform, url)

	var step models.IncomingExecutionStep
	err = json.NewDecoder(resp.Body).Decode(&step)
	if err != nil {
		log.Fatal(err)
		return bmodels.ExecutionSteps{}, err
	}

	return step.StepData, nil
}
