package alerts

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"

	log "github.com/sirupsen/logrus"
)

func SendAlert(cfg *config.Config, alert models.Alerts) {
	log.Info("Sending Alert")

	jsonPayload, err := json.Marshal(alert)
	if err != nil {
		log.Error(err)
		return
	}

	// Add authorization
	req, err := http.NewRequest("POST", cfg.Alertflow.URL+"/api/v1/alerts/", bytes.NewReader(jsonPayload))
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
		log.Error("Failed to send alert")
		return
	} else {
		log.Info("Alert Sent")
	}
}
