package payloads

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/pkg/models"
	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetData(cfg config.Config, payloadID string) (bmodels.Payloads, error) {

	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := cfg.Alertflow.URL + "/api/v1/payloads/" + payloadID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return bmodels.Payloads{}, err
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return bmodels.Payloads{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get payload from API: %s", url)
		err = fmt.Errorf("failed to get payload from API: %s", url)
		return bmodels.Payloads{}, err
	}

	log.Debugf("Payload data received from API: %s", url)

	var payload models.IncomingPayload
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		log.Fatal(err)
		return bmodels.Payloads{}, err
	}

	return payload.PayloadData, nil
}
