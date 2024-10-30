package actions

import "alertflow-runner/pkg/models"

func Init() []models.ActionDetails {
	var actions []models.ActionDetails

	CollectDataAction := CollectDataInit()
	CollectFlowDataAction := CollectFlowDataInit()
	CollectPayloadDataAction := CollectPayloadDataInit()
	PatternCheckAction := PatternCheckInit()
	FlowActionsCheckAction := FlowActionsCheckInit()
	logAction := LogInit()
	webhookAction := WebhookInit()
	WaitAction := WaitInit()
	PingAction := PingInit()
	PortAction := PortInit()
	InteractionAction := InteractionInit()

	actions = append(
		actions,
		CollectDataAction,
		CollectFlowDataAction,
		CollectPayloadDataAction,
		PatternCheckAction,
		FlowActionsCheckAction,
		logAction,
		webhookAction,
		WaitAction,
		PingAction,
		PortAction,
		InteractionAction,
	)
	return actions
}
