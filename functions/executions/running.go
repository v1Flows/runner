package executions

import (
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func SetToRunning(execution models.Execution) {
	execution.Running = true
	execution.Paused = false
	Running(execution)
}

func Running(execution models.Execution) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(execution)

	req, err := http.NewRequest("PUT", config.Config.Alertflow.URL+"/api/v1/executions/"+execution.ID.String(), payloadBuf)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error("Failed to update execution at API")
	}
}
