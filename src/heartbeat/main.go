package heartbeat

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func SendHeartbeat(api_url string, api_key string, runner_id string) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	url := api_url + "/api/runners/" + runner_id + "/heartbeat"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", api_key)
	for range time.Tick(time.Second * 10) {
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Fatalf("Failed to send heartbeat to API: %s", url)
		}
		log.Debugf("Heartbeat sent to API: %s", url)
	}
}
