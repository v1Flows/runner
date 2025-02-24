package alerts

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/AlertFlow/runner/config"
	"github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

func UpdateAlert(cfg config.Config, alert models.Alerts) {
	log.Info("Updating Alert")

	jsonPayload, err := json.Marshal(alert)
	if err != nil {
		log.Error(err)
		return
	}

	// Add authorization
	req, err := http.NewRequest("PUT", cfg.Alertflow.URL+"/api/v1/alerts/", bytes.NewReader(jsonPayload))
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
		log.Error("Failed to update alert")
		return
	} else {
		log.Info("Alert Updated")
	}
}
