package runner

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/platform"

	log "github.com/sirupsen/logrus"
)

func Busy(targetPlatform string, cfg config.Config, busy bool) {
	payload := models.Runners{
		ExecutingJob: busy,
	}

	url, apiKey, runnerID := platform.GetPlatformConfig(targetPlatform, cfg)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)
	req, err := http.NewRequest("PUT", url+"/api/v1/runners/"+runnerID+"/busy", payloadBuf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", apiKey)
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
		log.Error("Failed to set runner to busy at %s", targetPlatform)
		log.Error("Response: ", string(body))
	}
}
