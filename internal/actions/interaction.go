package actions

import (
	"encoding/json"
	"strconv"
	"time"

	"gitlab.justlab.xyz/alertflow-public/runner/pkg/executions"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

func InteractionInit() models.ActionDetails {
	params := []models.Param{
		{
			Key:         "Timeout",
			Type:        "number",
			Default:     0,
			Required:    false,
			Description: "Continue to the next step after the specified time (in seconds). 0 to disable",
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}

	return models.ActionDetails{
		ID:          "interaction",
		Name:        "Interaction",
		Description: "Wait for user interaction to continue",
		Icon:        "solar:hand-shake-linear",
		Type:        "interaction",
		Category:    "General",
		Function:    InteractionAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func InteractionAction(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	timeout := 0
	for _, param := range action.Params {
		if param.Key == "Timeout" {
			timeout, _ = strconv.Atoi(param.Value)
		}
	}

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:       step.ID,
		ActionID: action.ID.String(),
		ActionMessages: []string{
			`Waiting for user interaction`,
			`Timeout: ` + strconv.Itoa(timeout) + ` seconds`,
		},
		Pending:     false,
		Interactive: true,
		Running:     true,
		StartedAt:   time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	executions.SetToInteractionRequired(execution)

	var stepData models.ExecutionSteps

	// pull current action status from backend every 10 seconds
	startTime := time.Now()
	for {
		stepData, err = executions.GetStep(execution.ID.String(), step.ID.String())

		if stepData.Interacted {
			break
		} else {
			time.Sleep(5 * time.Second)
		}

		if err != nil {
			log.Error("Error getting step data: ", err)
			return nil, false, false, false, true
		}

		if timeout > 0 && time.Since(startTime).Seconds() >= float64(timeout) {
			log.Debug("Timeout reached while waiting for user interaction")
			err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
				ID: step.ID,
				ActionMessages: []string{
					"Interaction timed out",
					"Automatically approved & continuing to the next step",
				},
				Finished:            true,
				FinishedAt:          time.Now(),
				Interacted:          true,
				InteractionApproved: true,
				InteractionRejected: false,
			})
			if err != nil {
				return nil, false, false, false, true
			}

			stepData.Interacted = true
			stepData.InteractionApproved = true
			break
		}
	}

	executions.SetToRunning(execution)

	if stepData.InteractionRejected {
		err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID: step.ID,
			ActionMessages: []string{
				"Interaction rejected",
				"Execution canceled",
			},
			Running:             false,
			Canceled:            true,
			Finished:            true,
			FinishedAt:          time.Now(),
			Interacted:          true,
			InteractionRejected: true,
			InteractionApproved: false,
		})
		if err != nil {
			return nil, false, false, false, true
		}
		return nil, false, true, false, false
	} else if stepData.InteractionApproved {
		err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:                  step.ID,
			ActionMessages:      []string{"Interaction approved"},
			Running:             false,
			Finished:            true,
			FinishedAt:          time.Now(),
			Interacted:          true,
			InteractionRejected: false,
			InteractionApproved: true,
		})
		if err != nil {
			return nil, false, false, false, true
		}
		return nil, true, false, false, false
	}

	return nil, true, false, false, false
}
