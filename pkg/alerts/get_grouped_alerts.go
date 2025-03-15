package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetGroupedAlerts(cfg config.Config, flowID string, groupKeyIdentifier string) ([]bmodels.Alerts, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	request := bmodels.IncomingGroupedAlertsRequest{
		FlowID:                flowID,
		GroupAlertsIdentifier: groupKeyIdentifier,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(request)

	url := cfg.Alertflow.URL + "/api/v1/alerts/grouped"
	req, err := http.NewRequest("GET", url, payloadBuf)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return []bmodels.Alerts{}, err
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return []bmodels.Alerts{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get alerts from API: %s", url)
		err = fmt.Errorf("failed to get alerts from API: %s", url)
		return []bmodels.Alerts{}, err
	}

	log.Debugf("Alerts received from API: %s", url)

	var alerts models.IncomingAlerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	if err != nil {
		log.Fatal(err)
		return []bmodels.Alerts{}, err
	}

	return alerts.Alerts, nil
}
