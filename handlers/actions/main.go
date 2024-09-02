package handler_actions

import "alertflow-runner/models"

func Init() []models.ActionDetails {
	var actions []models.ActionDetails

	logAction := LogInit()
	webhookAction := WebhookInit()
	WaitAction := WaitInit()
	PingAction := PingInit()

	actions = append(actions, logAction, webhookAction, WaitAction, PingAction)
	return actions
}
