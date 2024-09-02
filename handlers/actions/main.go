package handler_actions

import "alertflow-runner/models"

func Init() []models.ActionDetails {
	var actions []models.ActionDetails

	logAction := LogInit()
	webhookAction := WebhookInit()
	WaitAction := WaitInit()
	PingAction := PingInit()
	PortAction := PortInit()

	actions = append(actions, logAction, webhookAction, WaitAction, PingAction, PortAction)
	return actions
}
