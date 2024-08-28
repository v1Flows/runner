package handler_actions

import (
	"alertflow-runner/functions/executions"
	"alertflow-runner/handlers/variables"
	"alertflow-runner/models"
	"time"

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

func LogAction() {
	log.WithFields(log.Fields{
		"Action":    variables.CurrentActionDetails.Name,
		"Type":      variables.CurrentActionDetails.Type,
		"Execution": variables.CurrentExecution.ID,
	}).Info("Log Action triggered")

	err := executions.UpdateStep(variables.CurrentExecution, models.ExecutionSteps{
		ID:             variables.CurrentActionStep.ID,
		ActionMessages: []string{"Log Action finished"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}
}
