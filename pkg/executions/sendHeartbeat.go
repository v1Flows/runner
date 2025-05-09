package executions

import (
	"net/http"

	"github.com/v1Flows/runner/config"
	"github.com/v1Flows/runner/pkg/platform"
	shared_models "github.com/v1Flows/shared-library/pkg/models"

	log "github.com/sirupsen/logrus"
)

func SendHeartbeat(cfg config.Config, execution shared_models.Executions, targetPlatform string) {
	url, apiKey := platform.GetPlatformConfigPlain(targetPlatform, cfg)

	req, err := http.NewRequest("PUT", url+"/api/v1/executions/"+execution.ID.String()+"/heartbeat", nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Authorization", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error("Failed to send execution heatbeat at " + targetPlatform + " API")
	}
}
