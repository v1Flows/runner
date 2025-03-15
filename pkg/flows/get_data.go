package flows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	bmodels "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/internal/common"
	"github.com/v1Flows/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func GetFlowData(platform string, cfg config.Config, flowID string) (bmodels.Flows, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url, apiKey, _ := common.GetPlatformConfig(platform, cfg)

	parsedUrl := url + "/api/v1/flows/" + flowID
	req, err := http.NewRequest("GET", parsedUrl, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return bmodels.Flows{}, err
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return bmodels.Flows{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get flow data from %s API: %s", platform, url)
		err = fmt.Errorf("failed to get flow data from %s API: %s", platform, url)
		return bmodels.Flows{}, err
	}

	log.Debugf("Flow data received from %s API: %s", platform, url)

	var flow models.IncomingFlow
	err = json.NewDecoder(resp.Body).Decode(&flow)
	if err != nil {
		log.Fatal(err)
		return bmodels.Flows{}, err
	}

	return flow.FlowData, nil
}
