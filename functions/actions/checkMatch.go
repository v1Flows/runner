package actions

import (
	"alertflow-runner/functions/executions"
	"alertflow-runner/models"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func checkMatch(execution models.Execution, action models.FlowActions, payload models.Payload, actionStepID string) (bool, error) {

	checkActionMatchDataStep, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Match Patterns",
		ActionMessages: []string{"Checking Action Match Patterns"},
		StartedAt:      time.Now(),
		ParentID:       actionStepID,
		IsHidden:       true,
	})
	if err != nil {
		log.Error("Error sending step:", err)
		return false, err
	}

	// convert payload to string
	payloadBytes, err := json.Marshal(payload.Payload)
	if err != nil {
		log.Error("Error converting payload to JSON:", err)
		return false, err
	}
	payloadString := string(payloadBytes)

	for _, matchPattern := range action.MatchPatterns {
		value := gjson.Get(payloadString, matchPattern.Key)

		if value.String() == matchPattern.Value {
			err := executions.UpdateStep(execution, models.ExecutionSteps{
				ID:             checkActionMatchDataStep.ID,
				ActionMessages: []string{`Match Pattern "` + matchPattern.Key + `": "` + matchPattern.Value + `" found`},
			})
			if err != nil {
				log.Error("Error updating step:", err)
				return false, err
			}
		} else {
			err = executions.UpdateStep(execution, models.ExecutionSteps{
				ID:             checkActionMatchDataStep.ID,
				ActionMessages: []string{`Match Pattern "` + matchPattern.Key + `": "` + matchPattern.Value + `" not found`},
				NoPatternMatch: true,
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				log.Error("Error updating step:", err)
				return false, err
			}
			return false, nil
		}
	}

	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             checkActionMatchDataStep.ID,
		ActionMessages: []string{"Match Patterns checked"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step:", err)
		return false, err
	}

	return true, nil
}
