package actions

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

func WebhookInit() models.ActionDetails {
	params := []models.Param{
		{
			Key:         "Url",
			Type:        "text",
			Required:    true,
			Description: "The URL to send the Webhook",
		},
		{
			Key:         "Method",
			Type:        "text",
			Default:     "POST",
			Required:    true,
			Description: "The HTTP Method to use",
		},
		{
			Key:         "Headers",
			Type:        "textarea",
			Required:    false,
			Description: "The headers to send",
		},
		{
			Key:         "Body",
			Type:        "textarea",
			Required:    false,
			Description: "The body to send",
		},
		{
			Key:         "Timeout",
			Type:        "number",
			Default:     10,
			Required:    true,
			Description: "The timeout in seconds",
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

func WebhookAction(execution models.Execution, step models.ExecutionSteps, action models.Actions) (finished bool, canceled bool, failed bool) {
	err := executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Webhook Action finished"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return true, false, false
}
