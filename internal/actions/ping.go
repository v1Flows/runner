package actions

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"encoding/json"
	"strconv"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	log "github.com/sirupsen/logrus"
)

func PingInit() models.ActionDetails {
	params := []models.Param{
		{
			Key:         "Target",
			Type:        "text",
			Default:     "www.alertflow.org",
			Required:    true,
			Description: "The target to ping",
		},
		{
			Key:         "Count",
			Type:        "number",
			Default:     3,
			Required:    false,
			Description: "Number of packets to send",
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}

	return models.ActionDetails{
		Name:        "Ping",
		Description: "Pings a target",
		Icon:        "solar:wi-fi-router-minimalistic-broken",
		Type:        "ping",
		Function:    PingAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func PingAction(execution models.Execution, step models.ExecutionSteps, action models.Actions) (finished bool, canceled bool, failed bool) {
	target := "www.alertflow.org"
	count := 3
	for _, param := range action.Params {
		if param.Key == "Target" {
			target = param.Value
		}
		if param.Key == "Count" {
			count, _ = strconv.Atoi(param.Value)
		}
	}

	err := executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{`Pinging: ` + target},
	})
	if err != nil {
		log.Error("Error updating step:", err)
	}

	pinger, err := probing.NewPinger(target)
	if err != nil {
		log.Error("Error creating pinger: ", err)
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"Error creating pinger: " + err.Error()},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error("Error updating step: ", err)
		}
		return false, false, true
	}
	pinger.Count = count
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Error("Error running pinger: ", err)
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"Error running pinger: " + err.Error()},
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error("Error updating step: ", err)
		}
		return false, false, true
	}

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID: step.ID,
		ActionMessages: []string{
			"Sent: " + strconv.Itoa(stats.PacketsSent),
			"Received: " + strconv.Itoa(stats.PacketsRecv),
			"Lost: " + strconv.Itoa(int(stats.PacketLoss)),
			"RTT min: " + stats.MinRtt.String(),
			"RTT max: " + stats.MaxRtt.String(),
			"RTT avg: " + stats.AvgRtt.String(),
		},
		Finished:   true,
		FinishedAt: time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Ping finished"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return true, false, false
}
