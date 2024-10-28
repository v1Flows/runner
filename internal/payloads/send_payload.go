package payloads

import (
	"alertflow-runner/config"
	"alertflow-runner/pkg/models"
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func SendPayload(payload models.Payload) {
	log.Info("Sending Payload")

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
		return
	}

	// Add authorization
	req, err := http.NewRequest("POST", config.Config.Alertflow.URL+"/api/v1/flows/"+payload.FlowID+"/payloads/", bytes.NewReader(jsonPayload))
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)

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
