package alerts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetData(cfg *config.Config, alertID string) (bmodels.Alerts, error) {

	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := cfg.Alertflow.URL + "/api/v1/alerts/" + alertID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return bmodels.Alerts{}, err
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return bmodels.Alerts{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get alert from API: %s", url)
		err = fmt.Errorf("failed to get alert from API: %s", url)
		return bmodels.Alerts{}, err
	}

	log.Debugf("Alert data received from API: %s", url)

	var alert models.IncomingAlert
	err = json.NewDecoder(resp.Body).Decode(&alert)
	if err != nil {
		log.Fatal(err)
		return bmodels.Alerts{}, err
	}

	return alert.AlertData, nil
}
