package executions

import (
	"alertflow-runner/handlers/busy"
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func EndWithError(execution models.Execution) {
	execution.FinishedAt = time.Now()
	execution.Running = false
	execution.Error = true
	End(execution)
}

func EndWithNoMatch(execution models.Execution) {
	execution.FinishedAt = time.Now()
	execution.Running = false
	execution.Error = false
	execution.NoMatch = true
	End(execution)
}

func EndWithGhost(execution models.Execution) {
	execution.FinishedAt = time.Now()
	execution.Running = false
	execution.Error = false
	execution.Ghost = true
	End(execution)
}

func EndSuccess(execution models.Execution) {
	execution.FinishedAt = time.Now()
	execution.Running = false
	execution.Error = false
	End(execution)
}

func End(execution models.Execution) {
	busy.Busy(false)

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
