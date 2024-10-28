package actions

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"encoding/json"
	"strconv"
	"time"

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
		Function:    InteractionAction,
		Params:      json.RawMessage(paramsJSON),
	}
}

func InteractionAction(execution models.Execution, step models.ExecutionSteps, action models.Actions) (finished bool, canceled bool, failed bool) {
	timeout := 0
	for _, param := range action.Params {
		if param.Key == "Timeout" {
			timeout, _ = strconv.Atoi(param.Value)
		}
	}

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{`Waiting for user interaction`},
		Interactive:    true,
	})
	if err != nil {
		log.Error("Error updating step:", err)
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
			return false, false, true
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
				log.Error("Error updating step: ", err)
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
			Finished:            true,
			FinishedAt:          time.Now(),
			Interacted:          true,
			InteractionRejected: true,
			InteractionApproved: false,
		})
		if err != nil {
			log.Error("Error updating step: ", err)
		}
		return false, true, false
	} else if stepData.InteractionApproved {
		err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:                  step.ID,
			ActionMessages:      []string{"Interaction approved"},
			Finished:            true,
			FinishedAt:          time.Now(),
			Interacted:          true,
			InteractionRejected: false,
			InteractionApproved: true,
		})
		if err != nil {
			log.Error("Error updating step: ", err)
		}
		return true, false, false
	}

	err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Interaction finished"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return true, false, false
}
