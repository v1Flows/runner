package runner

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func Busy(busy bool) {
	payload := models.Busy{
		ExecutingJob: busy,
	}

	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)
	req, err := http.NewRequest("PUT", cfg.Alertflow.URL+"/api/v1/runners/"+configManager.GetRunnerID()+"/busy", payloadBuf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 201 {
		log.Error("Failed to set runner to busy at AlertFlow")
		log.Error("Response: ", string(body))
	}
}
