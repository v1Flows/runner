package payloadhandler

import (
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"
)

func Init() []models.PayloadEndpoint {
	var payloadInjects []models.PayloadEndpoint

	// Alertmanager := AlertmanagerPayloadHandlerInit()

	payloadInjects = append(
		payloadInjects,
		// Alertmanager,
	)
	return payloadInjects
}
