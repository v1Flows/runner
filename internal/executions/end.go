package executions

import (
	"alertflow-runner/config"
	"alertflow-runner/internal/runner"
	"alertflow-runner/pkg/models"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func EndCanceled(execution models.Execution) {
	execution.FinishedAt = time.Now()
	execution.Running = false
	execution.Error = false
	execution.Canceled = true
	End(execution)
}

func EndNoPatternMatch(execution models.Execution) {
	execution.FinishedAt = time.Now()
	execution.Running = false
	execution.Error = false
	execution.NoPatternMatch = true
	End(execution)
}

func EndWithError(execution models.Execution) {
	execution.FinishedAt = time.Now()
	execution.Running = false
	execution.Error = true
	End(execution)
}

func EndSuccess(execution models.Execution) {
	execution.Running = false
	execution.Error = false
	execution.Finished = true
	execution.FinishedAt = time.Now()
	End(execution)
}

func End(execution models.Execution) {
	runner.Busy(false)

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
