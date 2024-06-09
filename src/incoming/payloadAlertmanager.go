package incoming

import (
	"alertflow-runner/src/config"
	"alertflow-runner/src/interactions/processing"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Payload struct {
	Payload json.RawMessage `json:"payload"`
}

type Receiver struct {
	Receiver string `json:"receiver"`
}

// AlertmanagerPayloadHandler processes the payload received from Alertmanager.
// It sends the payload to the Alertflow API.
func AlertmanagerPayloadHandler(ctx *gin.Context) {
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	payloadData := Payload{
		Payload: payload,
	}

	payloadBuf := new(bytes.Buffer)
	err = json.NewEncoder(payloadBuf).Encode(payloadData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode payload",
		})
		return
	}

	req, err := http.NewRequest("POST", config.Config.Alertflow.URL+"/api/payloads/", payloadBuf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create request",
		})
		return
	}

	req.Header.Set("Authorization", config.Config.Alertflow.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send request",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send payload to API",
		})
		return
	}

	var respPayload Payload
	err = json.NewDecoder(resp.Body).Decode(&respPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to decode response",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "payload sent",
		"response": respPayload,
	})

	receiverObj := new(Receiver)
	err = json.Unmarshal(respPayload.Payload, receiverObj)
	if err != nil {
		log.Error("Failed to unmarshal payload: ", err)
		return
	}
	processing.StartProcessing(receiverObj.Receiver)
}
