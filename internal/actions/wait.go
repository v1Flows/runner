package actions

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"encoding/json"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func WaitInit() models.ActionDetails {
	params := []models.Param{
		{
			Key:         "WaitTime",
			Type:        "number",
			Default:     10,
			Required:    true,
			Description: "The time to wait in seconds",
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}

	return models.ActionDetails{
		Name:        "Wait",
		Description: "Waits for a specified amount of time",
		Icon:        "solar:clock-circle-broken",
		Type:        "wait",
		Category:    "General",
		Function:    WaitAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func WaitAction(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	// get the waittime from the action params
	waitTime := 10
	for _, param := range action.Params {
		if param.Key == "WaitTime" {
			waitTime, _ = strconv.Atoi(param.Value)
		}
	}

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{`Waiting for ` + strconv.Itoa(waitTime) + ` seconds`},
		Pending:        false,
		Paused:         true,
		StartedAt:      time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	executions.SetToPaused(execution)

	time.Sleep(time.Duration(waitTime) * time.Second)

	executions.SetToRunning(execution)

	err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Wait Action finished"},
		Paused:         false,
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	return nil, true, false, false, false
}
