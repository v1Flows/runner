package actions

import "alertflow-runner/src/models"

func Init() []models.ActionDetails {
	var actions []models.ActionDetails

	logAction := LogInit()
	WebhookAction := WebhookInit()

	actions = append(actions, logAction, WebhookAction)
	return actions
}
