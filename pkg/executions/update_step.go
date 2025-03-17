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

func UpdateStep(cfg config.Config, executionID string, step shared_models.ExecutionSteps) error {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	targetPlatform, ok := platform.GetPlatformForExecution(executionID)
	if !ok {
		log.Error("Failed to get platform")
		return fmt.Errorf("failed to get platform")
	}

	url, apiKey, _ := platform.GetPlatformConfig(targetPlatform, cfg)

	req, err := http.NewRequest("PUT", url+"/api/v1/executions/"+executionID+"/steps/"+step.ID.String(), payloadBuf)
	if err != nil {
		log.Error(err)
		return err
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error("Failed to send execution step at " + targetPlatform + " API")
		return err
	}

	return nil
}
