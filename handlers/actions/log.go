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
		Description: "Prints a Log Message on Runner stdout",
		Icon:        "solar:clipboard-list-broken",
		Type:        "log",
		Function:    LogAction,
		Params:      nil,
	}
}

func LogAction(step models.ExecutionSteps, action models.Actions) bool {
	log.WithFields(log.Fields{
		"Execution": variables.CurrentExecution.ID,
		"StepID":    step.ID,
	}).Info("Log Action triggered")

	err := executions.UpdateStep(variables.CurrentExecution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Log Action finished"},
		Icon:           variables.CurrentActionDetails.Icon,
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return true
}
