package flow

import (
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetFlowData(execution models.Execution) (models.Flows, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := config.Config.Alertflow.URL + "/api/v1/flows/" + execution.FlowID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return models.Flows{}, err
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return models.Flows{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get flow data from API: %s", url)
		err = fmt.Errorf("failed to get flow data from API: %s", url)
		return models.Flows{}, err
	}

	log.Debugf("Flow data received from API: %s", url)

	var flow models.IncomingFlow
	err = json.NewDecoder(resp.Body).Decode(&flow)
	if err != nil {
		log.Fatal(err)
		return models.Flows{}, err
	}

	return flow.FlowData, nil
}
