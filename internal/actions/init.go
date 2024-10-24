package actions

import "alertflow-runner/pkg/models"

func Init() []models.ActionDetails {
	var actions []models.ActionDetails

	logAction := LogInit()
	webhookAction := WebhookInit()
	WaitAction := WaitInit()
	PingAction := PingInit()
	PortAction := PortInit()
	InteractionAction := InteractionInit()

	actions = append(actions, logAction, webhookAction, WaitAction, PingAction, PortAction, InteractionAction)
	return actions
}
