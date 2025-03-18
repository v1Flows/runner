package executions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/runner"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func EndCanceled(cfg config.Config, execution shared_models.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "canceled"
	End(cfg, execution)
}

func EndNoPatternMatch(cfg config.Config, execution shared_models.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "noPatternMatch"
	End(cfg, execution)
}

func EndWithError(cfg config.Config, execution shared_models.Executions) {
	execution.FinishedAt = time.Now()
	execution.Status = "error"
	End(cfg, execution)
}

func EndSuccess(cfg config.Config, execution shared_models.Executions) {
	execution.Status = "success"
	execution.FinishedAt = time.Now()
	End(cfg, execution)
}

func End(cfg config.Config, execution shared_models.Executions) {
	targetPlatform, ok := platform.GetPlatformForExecution(execution.ID.String())
	if !ok {
		log.Error("Failed to get platform")
		return
	}

	url, apiKey := platform.GetPlatformConfigPlain(targetPlatform, cfg)

	runner.Busy(targetPlatform, cfg, false)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(execution)

	req, err := http.NewRequest("PUT", url+"/api/v1/executions/"+execution.ID.String(), payloadBuf)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error("Failed to update execution at " + targetPlatform + " API")
	}
}
