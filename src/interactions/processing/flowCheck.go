package processing

import (
	"alertflow-runner/src/config"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func CheckForFlow(flowID string) bool {
	req, err := http.NewRequest("GET", config.Config.Alertflow.URL+"/api/flows/"+flowID, nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
	}

	if resp.StatusCode != 200 {
		log.Error("Failed to check for flow: ", flowID)
		return false
	}

	defer resp.Body.Close()

	return true
}
