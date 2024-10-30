package common

import (
	"alertflow-runner/pkg/models"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartWorker(api_url string, api_key string, runner_id string) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	url := api_url + "/api/v1/executions/" + runner_id + "/pending"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", api_key)
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}

		if resp.StatusCode != 200 {
			log.Errorf("Failed to get waiting executions from API: %s", url)
			panic("Failed to get waiting executions from API")
		}

		log.Debugf("Executions received from API: %s", url)

		var executions models.Executions
		err = json.NewDecoder(resp.Body).Decode(&executions)
		if err != nil {
			log.Fatal(err)
		}

		for _, execution := range executions.Executions {
			go startProcessing(execution)
		}
	}
}
