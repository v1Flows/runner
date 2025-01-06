package executions

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func UpdateStep(executionID string, step models.ExecutionSteps) error {

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(step)

	req, err := http.NewRequest("PUT", config.Config.Alertflow.URL+"/api/v1/executions/"+executionID+"/steps/"+step.ID.String(), payloadBuf)
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
