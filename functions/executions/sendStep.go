package executions

import (
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func SendStep(execution models.Execution, step models.ExecutionSteps) {

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	req, err := http.NewRequest("POST", config.Config.Alertflow.URL+"/api/executions/"+execution.ID.String()+"/steps", payloadBuf)
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
}
