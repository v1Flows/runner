package actions

import (
	"alertflow-runner/models"
	"encoding/json"
)

func WebhookInit() models.ActionDetails {
	params := []models.Param{
		{
			Key:      "url",
			Type:     "string",
			Required: true,
		},
		{
			Key:      "method",
			Type:     "string",
			Default:  "POST",
			Required: false,
		},
		{
			Key:      "headers",
			Type:     "object",
			Required: false,
		},
		{
			Key:      "body",
			Type:     "string",
			Required: false,
		},
		{
			Key:      "timeout",
			Type:     "string",
			Default:  "10s",
			Required: false,
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		// handle error
	}

	return models.ActionDetails{
		Name:        "Webhook",
		Description: "Sends an HTTP Webhook",
		Type:        "webhook",
		Params:      json.RawMessage(paramsJSON),
	}
}
