package register

import (
	"alertflow-runner/src/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Register struct {
	Registered                bool            `json:"registered"`
	AvailableActions          json.RawMessage `json:"available_actions"`
	AvailablePayloadInjectors json.RawMessage `json:"available_payload_injectors"`
	LastHeartbeat             time.Time       `json:"last_heartbeat"`
	RunnerVersion             string          `json:"runner_version"`
}

func RegisterAtAPI(api_url string, api_key string, runner_id string, version string, actions []models.ActionDetails) {
	register := Register{
		Registered:                true,
		LastHeartbeat:             time.Now(),
		RunnerVersion:             version,
		AvailablePayloadInjectors: json.RawMessage(`[]`),
	}

	// Convert actions to JSON
	actionsJSON, err := json.Marshal(actions)
	if err != nil {
		log.Fatal(err)
	}

	register.AvailableActions = json.RawMessage(actionsJSON)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(register)
	req, err := http.NewRequest("POST", api_url+"/api/runners/"+runner_id+"/register", payloadBuf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", api_key)
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
		log.Error("Failed to register at API: ", api_url+"/api/runners/"+runner_id+"/register")
		log.Error("Response: ", string(body))
		panic("Failed to register at API")
	}
	log.Info("Registered at API: ", api_url+"/api/runners/"+runner_id+"/register")
}
