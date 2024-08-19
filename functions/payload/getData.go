package payload

import (
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetData(execution models.Execution) (models.Payload, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := config.Config.Alertflow.URL + "/api/payloads/" + execution.PayloadID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return models.Payload{}, err
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return models.Payload{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get payload from API: %s", url)
		err = fmt.Errorf("failed to get payload from API: %s", url)
		return models.Payload{}, err
	}

	log.Debugf("Payload data received from API: %s", url)

	var payload models.IncomingPayload
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		log.Fatal(err)
		return models.Payload{}, err
	}

	return payload.PayloadData, nil
}
