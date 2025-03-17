package executions

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func SendStep(cfg config.Config, execution shared_models.Executions, step shared_models.ExecutionSteps) (shared_models.ExecutionSteps, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	platform, ok := platform.GetPlatformForExecution(execution.ID.String())
	if !ok {
		log.Error("Failed to get platform")
		return shared_models.ExecutionSteps{}, nil
	}

	url, apiKey, _ := common.GetPlatformConfig(platform, cfg)

	req, err := http.NewRequest("POST", url+"/api/v1/executions/"+execution.ID.String()+"/steps", payloadBuf)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		log.Error("Failed to send execution step at " + platform + " API")
	}

	var stepResponse shared_models.ExecutionSteps
	err = json.NewDecoder(resp.Body).Decode(&stepResponse)
	if err != nil {
		log.Error(err)
		return shared_models.ExecutionSteps{}, err
	}

	return stepResponse, nil
}
