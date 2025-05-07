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

func EndCanceled(cfg config.Config, execution shared_models.Executions, targetPlatform string) {
	execution.FinishedAt = time.Now()
	execution.Status = "canceled"
	End(cfg, execution, targetPlatform)
}

func EndNoPatternMatch(cfg config.Config, execution shared_models.Executions, targetPlatform string) {
	execution.FinishedAt = time.Now()
	execution.Status = "noPatternMatch"
	End(cfg, execution, targetPlatform)
}

func EndWithError(cfg config.Config, execution shared_models.Executions, targetPlatform string) {
	execution.FinishedAt = time.Now()
	execution.Status = "error"
	End(cfg, execution, targetPlatform)
}

func EndWithRecovered(cfg config.Config, execution shared_models.Executions, targetPlatform string) {
	execution.FinishedAt = time.Now()
	execution.Status = "recovered"
	End(cfg, execution, targetPlatform)
}

func EndSuccess(cfg config.Config, execution shared_models.Executions, targetPlatform string) {
	execution.Status = "success"
	execution.FinishedAt = time.Now()
	End(cfg, execution, targetPlatform)
}

func End(cfg config.Config, execution shared_models.Executions, targetPlatform string) {
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
