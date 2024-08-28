package handler_actions

import "alertflow-runner/models"

func Init() []models.ActionDetails {
	var actions []models.ActionDetails

	logAction := LogInit()
	WebhookAction := WebhookInit()

	actions = append(actions, logAction, WebhookAction)
	return actions
}
