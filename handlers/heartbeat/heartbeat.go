package heartbeat

import (
	"alertflow-runner/handlers/config"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func SendHeartbeat() {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	url := config.Config.Alertflow.URL + "/api/runners/" + config.Config.RunnerID + "/heartbeat"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Errorf("Failed to send heartbeat to AlertFlow")
			log.Errorf("Response: %s", body)
		}
		log.Debugf("Heartbeat sent to AlertFlow")
	}
}
