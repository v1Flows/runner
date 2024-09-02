package handler_actions

import (
	"alertflow-runner/functions/executions"
	"alertflow-runner/handlers/variables"
	"alertflow-runner/models"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

func WebhookInit() models.ActionDetails {
	params := []models.Param{
		{
			Key:      "Url",
			Type:     "text",
			Required: true,
		},
		{
			Key:      "Method",
			Type:     "text",
			Default:  "POST",
			Required: true,
		},
		{
			Key:      "Headers",
			Type:     "textarea",
			Required: false,
		},
		{
			Key:      "Body",
			Type:     "textarea",
			Required: false,
		},
		{
			Key:      "Timeout",
			Type:     "number",
			Default:  10,
			Required: true,
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}

	return models.ActionDetails{
		Name:        "Webhook",
		Description: "Sends an HTTP Webhook",
		Icon:        "solar:global-broken",
		Type:        "webhook",
		Function:    WebhookAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func WebhookAction(step models.ExecutionSteps, action models.Actions) bool {
	err := executions.UpdateStep(variables.CurrentExecution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Webhook Action finished"},
		Icon:           variables.CurrentActionDetails.Icon,
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return true
}
