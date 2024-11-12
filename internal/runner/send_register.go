package runner

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gitlab.justlab.xyz/alertflow-public/runner/config"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func RegisterAtAPI(version string, actions []models.ActionDetails, payloadInjectors []models.PayloadEndpoint) {
	register := models.Register{
		Registered:    true,
		LastHeartbeat: time.Now(),
		Version:       version,
		Mode:          config.Config.Mode,
	}

	// Convert actions to JSON
	actionsJSON, err := json.Marshal(actions)
	if err != nil {
		log.Fatal(err)
	}
	register.AvailableActions = json.RawMessage(actionsJSON)

	// Convert payloadInjectors to JSON
	payloadInjectorsJSON, err := json.Marshal(payloadInjectors)
	if err != nil {
		log.Fatal(err)
	}
	register.AvailablePayloadInjectors = json.RawMessage(payloadInjectorsJSON)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(register)
	req, err := http.NewRequest("PUT", config.Config.Alertflow.URL+"/api/v1/runners/"+config.Config.RunnerID+"/register", payloadBuf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
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
		log.Error("Failed to register at AlertFlow")
		log.Error("Response: ", string(body))
		panic("Failed to register at AlertFlow")
	}
	log.Info("Runner registered at AlertFlow")
}
