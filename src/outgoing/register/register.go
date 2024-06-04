package register

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Register struct {
	Registered       bool            `json:"registered"`
	AvailableActions json.RawMessage `json:"available_actions"`
	LastHeartbeat    sql.NullTime    `json:"last_heartbeat"`
	RunnerVersion    string          `json:"runner_version"`
}

func RegisterAtAPI(api_url string, api_key string, runner_id string, version string) {
	register := Register{
		Registered: true,
		AvailableActions: json.RawMessage(`[{
			"name": "log",
			"type": "action",
			"description": "Post Log Message on API Backend Server",
			"params": {}
		}]`),
		LastHeartbeat: sql.NullTime{Time: time.Now(), Valid: true},
		RunnerVersion: version,
	}

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

	if resp.StatusCode != 201 {
		log.Fatal("Failed to register at API: ", api_url+"/api/runners/"+runner_id+"/register")
	}
	log.Info("Registered at API: ", api_url+"/api/runners/"+runner_id+"/register")
}
