package flows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/models"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetFlowData(cfg config.Config, flowID string, targetPlatform string) (shared_models.Flows, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url, apiKey, _ := platform.GetPlatformConfig(targetPlatform, cfg)

	parsedUrl := url + "/api/v1/flows/" + flowID
	req, err := http.NewRequest("GET", parsedUrl, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return shared_models.Flows{}, err
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return shared_models.Flows{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get flow data from %s API: %s", targetPlatform, url)
		err = fmt.Errorf("failed to get flow data from %s API: %s", targetPlatform, url)
		return shared_models.Flows{}, err
	}

	log.Debugf("Flow data received from %s API: %s", targetPlatform, url)

	var flow models.IncomingFlow
	err = json.NewDecoder(resp.Body).Decode(&flow)
	if err != nil {
		log.Fatal(err)
		return shared_models.Flows{}, err
	}

	return flow.FlowData, nil
}
