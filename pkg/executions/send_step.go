package executions

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/AlertFlow/runner/config"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

func SendStep(execution bmodels.Executions, step bmodels.ExecutionSteps) (bmodels.ExecutionSteps, error) {
	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	req, err := http.NewRequest("POST", cfg.Alertflow.URL+"/api/v1/executions/"+execution.ID.String()+"/steps", payloadBuf)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		log.Error("Failed to send execution step at API")
	}

	var stepResponse bmodels.ExecutionSteps
	err = json.NewDecoder(resp.Body).Decode(&stepResponse)
	if err != nil {
		log.Error(err)
		return bmodels.ExecutionSteps{}, err
	}

	return stepResponse, nil
}
