package payloads

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/AlertFlow/runner/config"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

func SendPayload(payload bmodels.Payloads) {
	log.Info("Sending Payload")

	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
		return
	}

	// Add authorization
	req, err := http.NewRequest("POST", cfg.Alertflow.URL+"/api/v1/flows/"+payload.FlowID+"/payloads/", bytes.NewReader(jsonPayload))
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		log.Error("Failed to send payload")
		return
	} else {
		log.Info("Payload Sent")
	}
}
