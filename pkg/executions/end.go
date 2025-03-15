package executions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/runner"

	log "github.com/sirupsen/logrus"
)

func EndCanceled(cfg config.Config, execution bmodels.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "canceled"
	End(cfg, execution)
}

func EndNoPatternMatch(cfg config.Config, execution bmodels.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "noPatternMatch"
	End(cfg, execution)
}

func EndWithError(cfg config.Config, execution bmodels.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "error"
	End(cfg, execution)
}

func EndSuccess(cfg config.Config, execution bmodels.Executions) {
	execution.Status = "success"
	execution.FinishedAt = time.Now()
	End(cfg, execution)
}

func End(cfg config.Config, execution bmodels.Executions) {
	runner.Busy(cfg, false)

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
