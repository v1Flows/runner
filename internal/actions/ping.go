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
		ID:          "ping",
		Name:        "Ping",
		Description: "Pings a target",
		Icon:        "solar:wi-fi-router-minimalistic-broken",
		Type:        "ping",
		Category:    "Network",
		Function:    PingAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func PingAction(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
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

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{`Pinging: ` + target},
		Pending:        false,
		StartedAt:      time.Now(),
		Running:        true,
	})
	if err != nil {
		return nil, false, false, false, true
	}

	pinger, err := probing.NewPinger(target)
	if err != nil {
		log.Error("Error creating pinger: ", err)
		err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"Error creating pinger: " + err.Error()},
			Running:        false,
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return nil, false, false, false, true
		}
		return nil, false, false, false, true
	}
	pinger.Count = count
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Error("Error running pinger: ", err)
		err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"Error running pinger: " + err.Error()},
			Running:        false,
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			return nil, false, false, false, true
		}
		return nil, false, false, false, true
	}

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID: step.ID,
		ActionMessages: []string{
			"Sent: " + strconv.Itoa(stats.PacketsSent),
			"Received: " + strconv.Itoa(stats.PacketsRecv),
			"Lost: " + strconv.Itoa(int(stats.PacketLoss)),
			"RTT min: " + stats.MinRtt.String(),
			"RTT max: " + stats.MaxRtt.String(),
			"RTT avg: " + stats.AvgRtt.String(),
			"Ping finished",
		},
		Running:    false,
		Finished:   true,
		FinishedAt: time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	return nil, true, false, false, false
}
