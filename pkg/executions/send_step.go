package executions

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func SendStep(execution models.Execution, step models.ExecutionSteps) (models.ExecutionSteps, error) {

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	req, err := http.NewRequest("POST", config.Config.Alertflow.URL+"/api/v1/executions/"+execution.ID.String()+"/steps", payloadBuf)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		log.Error("Failed to send execution step at API")
	}

	var stepResponse models.ExecutionSteps
	err = json.NewDecoder(resp.Body).Decode(&stepResponse)
	if err != nil {
		log.Error(err)
		return models.ExecutionSteps{}, err
	}

	return stepResponse, nil
}
