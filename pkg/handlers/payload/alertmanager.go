package payloadhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"gitlab.justlab.xyz/alertflow-public/runner/config"
	internal_payloads "gitlab.justlab.xyz/alertflow-public/runner/internal/payloads"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Receiver struct {
	Receiver string `json:"receiver"`
}

func AlertmanagerPayloadHandlerInit() models.PayloadInjector {
	return models.PayloadInjector{
		Name:     "Alertmanager",
		Type:     "alertmanager",
		Endpoint: "/alertmanager",
	}
}

func AlertmanagerPayloadHandler(context *gin.Context) {
	log.Info("Received Alertmanager Payload")
	incPayload, err := io.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// get receiver from payload
	receiver := Receiver{}
	json.Unmarshal(incPayload, &receiver)

	payloadData := models.Payload{
		Payload:  incPayload,
		FlowID:   receiver.Receiver,
		RunnerID: config.Config.RunnerID,
		Endpoint: "alertmanager",
	}

	internal_payloads.SendPayload(payloadData)
}
