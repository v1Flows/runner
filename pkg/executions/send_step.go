package executions

import (
	"bytes"
	"encoding/json"
	"net/http"

	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"

	log "github.com/sirupsen/logrus"
)

func SendStep(cfg config.Config, execution bmodels.Executions, step bmodels.ExecutionSteps) (bmodels.ExecutionSteps, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	platform, ok := GetPlatformForExecution(execution.ID.String())
	if !ok {
		log.Error("Failed to get platform")
		return bmodels.ExecutionSteps{}, nil
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
		log.Error("Failed to send execution step at %s API", platform)
	}

	var stepResponse bmodels.ExecutionSteps
	err = json.NewDecoder(resp.Body).Decode(&stepResponse)
	if err != nil {
		log.Error(err)
		return bmodels.ExecutionSteps{}, err
	}

	return stepResponse, nil
}
