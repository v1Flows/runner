package handler_actions

import (
	"alertflow-runner/models"

	log "github.com/sirupsen/logrus"
)

func LogInit() models.ActionDetails {
	return models.ActionDetails{
		Name:        "Log Message",
		Description: "Prints an Log Message on Runner stdout",
		Icon:        "solar:clipboard-list-broken",
		Type:        "log",
		Function:    LogAction,
		Params:      nil,
	}
}

func LogAction(subActionStepID string, execution models.Execution, action models.FlowActions, payload models.Payload) {
	log.Info("Log Action triggered")
}
