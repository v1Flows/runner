package executions

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func SetToInteractionRequired(execution models.Execution) {
	execution.Running = false
	execution.Paused = false
	execution.Pending = false
	execution.Canceled = false
	execution.InteractionRequired = true
	InteractionRequired(execution)
}

func InteractionRequired(execution models.Execution) {
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
