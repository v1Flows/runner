package runner

import (
	"io"
	"net/http"
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"

	log "github.com/sirupsen/logrus"
)

func SendHeartbeat(platform string) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	configManager := config.GetInstance()
	cfg := configManager.GetConfig()

	url, apiKey, runnerID := common.GetPlatformConfig(platform, cfg)

	parsedUrl := url + "/api/v1/runners/" + runnerID + "/heartbeat"
	req, err := http.NewRequest("PUT", parsedUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", apiKey)
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		var resp *http.Response
		var err error

		for i := 0; i < 3; i++ {
			resp, err = client.Do(req)
			if err != nil {
				log.Errorf("Failed to send request: %v", err)
				time.Sleep(5 * time.Second) // Add delay before retrying
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("Failed to read response body: %v", err)
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			if resp.StatusCode == 200 {
				log.Debugf("Heartbeat sent to %s", platform)
				break
			} else {
				log.Errorf("Failed to send heartbeat to %s, attempt %d", platform, i+1)
				log.Errorf("Response: %s", body)
				time.Sleep(5 * time.Second) // Add delay before retrying
			}
		}
		if resp.StatusCode != 200 {
			log.Fatalf("Failed to send heartbeat to %s after 3 attempts", platform)
		}
	}
}
