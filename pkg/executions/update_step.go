package executions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"

	log "github.com/sirupsen/logrus"
)

func UpdateStep(cfg config.Config, executionID string, step bmodels.ExecutionSteps) error {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	platform, ok := GetPlatformForExecution(executionID)
	if !ok {
		log.Error("Failed to get platform")
		return fmt.Errorf("Failed to get platform")
	}

	url, apiKey, _ := common.GetPlatformConfig(platform, cfg)

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
		log.Error("Failed to send execution step at %s API", platform)
		return err
	}

	return nil
}
