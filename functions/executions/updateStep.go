package executions

import (
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func UpdateStep(execution models.Execution, step models.ExecutionSteps) error {

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	req, err := http.NewRequest("PUT", config.Config.Alertflow.URL+"/api/v1/executions/"+execution.ID.String()+"/steps/"+step.ID.String(), payloadBuf)
	if err != nil {
		log.Error(err)
		return err
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error("Failed to send execution step at API")
		return err
	}

	return nil
}
