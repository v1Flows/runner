package executions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetExecutionByID(cfg config.Config, executionID string, targetPlatform string) (shared_models.Executions, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url, apiKey := platform.GetPlatformConfigPlain(targetPlatform, cfg)

	parsedUrl := url + "/api/v1/executions/" + executionID
	req, err := http.NewRequest("GET", parsedUrl, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return shared_models.Executions{}, err
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return shared_models.Executions{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get execution data from %s API: %s", targetPlatform, url)
		err = fmt.Errorf("failed to get execution data from %s API: %s", targetPlatform, url)
		return shared_models.Executions{}, err
	}

	log.Debugf("Step data received from %s API: %s", targetPlatform, url)

	var execution shared_models.Executions
	err = json.NewDecoder(resp.Body).Decode(&execution)
	if err != nil {
		log.Fatal(err)
		return shared_models.Executions{}, err
	}

	return execution, nil
}
