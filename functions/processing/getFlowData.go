package processing

import (
	"alertflow-runner/handlers/config"
	"alertflow-runner/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetFlowData(execution models.Execution) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := config.Config.Alertflow.URL + "/api/flows/" + execution.FlowID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get waiting executions from API: %s", url)
		panic("Failed to get waiting executions from API")
	}

	log.Debugf("Flow data received from API: %s", url)

	var flow models.Flows
	err = json.NewDecoder(resp.Body).Decode(&flow)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(flow)
}
