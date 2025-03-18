package flows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	af_models "github.com/v1Flows/alertFlow/services/backend/pkg/models"
	ef_models "github.com/v1Flows/exFlow/services/backend/pkg/models"
	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/models"
	"github.com/v1Flows/runner/pkg/platform"

	log "github.com/sirupsen/logrus"
)

func GetFlowData(cfg config.Config, flowID string, targetPlatform string) (exFlow ef_models.Flows, alertFlow af_models.Flows, err error) {
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url, apiKey := platform.GetPlatformConfigPlain(targetPlatform, cfg)

	parsedUrl := url + "/api/v1/flows/" + flowID
	req, err := http.NewRequest("GET", parsedUrl, nil)
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return ef_models.Flows{}, af_models.Flows{}, err
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return ef_models.Flows{}, af_models.Flows{}, err
	}

	if resp.StatusCode != 200 {
		log.Errorf("Failed to get flow data from %s API: %s", targetPlatform, url)
		err = fmt.Errorf("failed to get flow data from %s API: %s", targetPlatform, url)
		return ef_models.Flows{}, af_models.Flows{}, err
	}

	log.Debugf("Flow data received from %s API: %s", targetPlatform, url)

	if targetPlatform == "alertflow" {
		var flow models.IncomingAfFlow
		err := json.NewDecoder(resp.Body).Decode(&flow)
		if err != nil {
			log.Fatal(err)
			return ef_models.Flows{}, af_models.Flows{}, err
		}
		return ef_models.Flows{}, flow.FlowData, nil
	} else if targetPlatform == "exflow" {
		var flow models.IncomingEfFlow
		err := json.NewDecoder(resp.Body).Decode(&flow)
		if err != nil {
			log.Fatal(err)
			return ef_models.Flows{}, af_models.Flows{}, err
		}
		return flow.FlowData, af_models.Flows{}, nil
	}

	return ef_models.Flows{}, af_models.Flows{}, fmt.Errorf("unknown target platform: %s", targetPlatform)
}
