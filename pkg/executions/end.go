package executions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/internal/runner"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

func EndCanceled(execution bmodels.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "canceled"
	End(execution)
}

func EndNoPatternMatch(execution bmodels.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "noPatternMatch"
	End(execution)
}

func EndWithError(execution bmodels.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "error"
	End(execution)
}

func EndSuccess(execution bmodels.Executions) {
	execution.Status = "success"
	execution.FinishedAt = time.Now()
	End(execution)
}

func End(execution bmodels.Executions) {
	runner.Busy(false)

	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(execution)

	req, err := http.NewRequest("PUT", cfg.Alertflow.URL+"/api/v1/executions/"+execution.ID.String(), payloadBuf)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error("Failed to update execution at API")
	}
}
