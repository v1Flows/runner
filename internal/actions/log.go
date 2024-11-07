package actions

import (
	"time"

	"gitlab.justlab.xyz/alertflow-public/runner/internal/executions"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func LogInit() models.ActionDetails {
	return models.ActionDetails{
		ID:          "log",
		Name:        "Log Message",
		Description: "Prints a Log Message on Runner stdout",
		Icon:        "solar:clipboard-list-broken",
		Type:        "log",
		Category:    "Utility",
		Function:    LogAction,
		Params:      nil,
	}
}

func LogAction(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	log.WithFields(log.Fields{
		"Execution": execution.ID,
		"StepID":    step.ID,
	}).Info("Log Action triggered")

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{"Log Action finished"},
		Pending:        false,
		Finished:       true,
		StartedAt:      time.Now(),
		FinishedAt:     time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	return nil, true, false, false, false
}
