package processing

import (
	log "github.com/sirupsen/logrus"
)

func StartProcessing(flowID string) {
	flow_check := CheckForFlow(flowID)
	if !flow_check {
		log.Info("Flow not found, exiting")
		return
	}
}
