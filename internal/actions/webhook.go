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
		ID:          "webhook",
		Name:        "Webhook",
		Description: "Sends an HTTP Webhook",
		Icon:        "solar:global-broken",
		Type:        "webhook",
		Category:    "Network",
		Function:    WebhookAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func WebhookAction(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{"Webhook Action finished"},
		Pending:        false,
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	return nil, true, false, false, false
}
