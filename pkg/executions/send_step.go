package executions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func SendStep(cfg config.Config, execution shared_models.Executions, step shared_models.ExecutionSteps, targetPlatform string) (shared_models.ExecutionSteps, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	url, apiKey := platform.GetPlatformConfigPlain(targetPlatform, cfg)

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
		log.Error("Failed to send execution step at " + targetPlatform + " API")
		return shared_models.ExecutionSteps{}, fmt.Errorf("failed to send execution step at " + targetPlatform + " api")
	}

	var stepResponse shared_models.ExecutionSteps
	err = json.NewDecoder(resp.Body).Decode(&stepResponse)
	if err != nil {
		log.Error(err)
		return shared_models.ExecutionSteps{}, err
	}

	return stepResponse, nil
}
