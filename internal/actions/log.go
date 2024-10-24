package actions

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
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

func LogAction(execution models.Execution, step models.ExecutionSteps, action models.Actions) (finished bool, canceled bool, failed bool) {
	log.WithFields(log.Fields{
		"Execution": execution.ID,
		"StepID":    step.ID,
	}).Info("Log Action triggered")

	err := executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Log Action finished"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return true, false, false
}
