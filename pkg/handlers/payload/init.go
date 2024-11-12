package payloadhandler

import (
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"
)

func Init() []models.PayloadInjector {
	var payloadInjects []models.PayloadInjector

	// Alertmanager := AlertmanagerPayloadHandlerInit()

	payloadInjects = append(
		payloadInjects,
		// Alertmanager,
	)
	return payloadInjects
}
