package payloadhandler

import (
	"alertflow-runner/pkg/models"
)

func Init() []models.PayloadInjector {
	var payloadInjects []models.PayloadInjector

	Alertmanager := AlertmanagerPayloadHandlerInit()

	payloadInjects = append(
		payloadInjects,
		Alertmanager,
	)
	return payloadInjects
}
