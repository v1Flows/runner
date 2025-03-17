package executions

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/platform"

	log "github.com/sirupsen/logrus"
)

func SetToInteractionRequired(cfg config.Config, execution models.Executions) {
	execution.Status = "interactionWaiting"
	InteractionWaiting(cfg, execution)
}

func InteractionWaiting(cfg config.Config, execution models.Executions) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(execution)

	platform, ok := platform.GetPlatformForExecution(execution.ID.String())
	if !ok {
		log.Error("Failed to get platform")
	}

	url, apiKey, _ := common.GetPlatformConfig(platform, cfg)

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
		log.Error("Failed to update execution at %s API", platform)
	}
}
