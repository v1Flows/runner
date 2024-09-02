package handler_actions

import (
	"alertflow-runner/functions/executions"
	"alertflow-runner/handlers/variables"
	"alertflow-runner/models"
	"encoding/json"
	"net"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func PortInit() models.ActionDetails {
	params := []models.Param{
		{
			Key:         "Host",
			Type:        "text",
			Default:     "myhost",
			Required:    true,
			Description: "The host to check for the port",
		},
		{
			Key:         "Port",
			Type:        "number",
			Default:     22,
			Required:    true,
			Description: "The port to check",
		},
		{
			Key:         "Timeout",
			Type:        "number",
			Default:     3,
			Required:    false,
			Description: "Timeout in seconds",
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}

	return models.ActionDetails{
		Name:        "Port Check",
		Description: "Checks if a port is open",
		Icon:        "solar:wi-fi-router-broken",
		Type:        "port",
		Function:    PortAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func PortAction(step models.ExecutionSteps, action models.Actions) bool {
	host := "myhost"
	port := 22
	timeout := 3
	for _, param := range action.Params {
		if param.Key == "Host" {
			host = param.Value
		}
		if param.Key == "Port" {
			port, _ = strconv.Atoi(param.Value)
		}
		if param.Key == "Timeout" {
			timeout, _ = strconv.Atoi(param.Value)
		}
	}

	err := executions.UpdateStep(variables.CurrentExecution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Checking port " + strconv.Itoa(port) + " on " + host},
	})
	if err != nil {
		log.Error("Error updating step:", err)
	}

	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		err = executions.UpdateStep(variables.CurrentExecution, models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"Port is closed"},
			Icon:           variables.CurrentActionDetails.Icon,
			Error:          true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error("Error updating step: ", err)
		}
		return false
	} else {
		if conn != nil {
			err = executions.UpdateStep(variables.CurrentExecution, models.ExecutionSteps{
				ID:             step.ID,
				ActionMessages: []string{"Port is open"},
				Icon:           variables.CurrentActionDetails.Icon,
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				log.Error("Error updating step: ", err)
			}
			defer conn.Close()
		} else {
			err = executions.UpdateStep(variables.CurrentExecution, models.ExecutionSteps{
				ID:             step.ID,
				ActionMessages: []string{"Port is closed"},
				Icon:           variables.CurrentActionDetails.Icon,
				Error:          true,
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				log.Error("Error updating step: ", err)
			}
			return false
		}
	}

	err = executions.UpdateStep(variables.CurrentExecution, models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Port check finished"},
		Icon:           variables.CurrentActionDetails.Icon,
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return true
}
