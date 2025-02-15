package flows

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

func GetFlowData(cfg config.Config, flowID string) (bmodels.Flows, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := cfg.Alertflow.URL + "/api/v1/flows/" + flowID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return bmodels.Flows{}, err
	}
	req.Header.Set("Authorization", cfg.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return bmodels.Flows{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get flow data from API: %s", url)
		err = fmt.Errorf("failed to get flow data from API: %s", url)
		return bmodels.Flows{}, err
	}

	log.Debugf("Flow data received from API: %s", url)

	var flow models.IncomingFlow
	err = json.NewDecoder(resp.Body).Decode(&flow)
	if err != nil {
		log.Fatal(err)
		return bmodels.Flows{}, err
	}

	return flow.FlowData, nil
}
