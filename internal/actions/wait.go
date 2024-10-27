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
		Function:    WaitAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func WaitAction(execution models.Execution, step models.ExecutionSteps, action models.Actions) (finished bool, canceled bool, failed bool) {
	// get the waittime from the action params
	waitTime := 10
	for _, param := range action.Params {
		if param.Key == "WaitTime" {
			waitTime, _ = strconv.Atoi(param.Value)
		}
	}

	err := executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{`Waiting for: ` + strconv.Itoa(waitTime) + ` seconds`},
	})
	if err != nil {
		log.Error("Error updating step:", err)
	}

	executions.SetToPaused(execution)

	time.Sleep(time.Duration(waitTime) * time.Second)

	executions.SetToRunning(execution)

	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Wait Action finished"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return true, false, false
}
