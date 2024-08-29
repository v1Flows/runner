package actions

import (
	"alertflow-runner/functions/executions"
	"alertflow-runner/models"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func checkExclude(execution models.Execution, action models.FlowActions, payload models.Payload, actionStepID string) (bool, error) {
	checkActionExcludeDataStep, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Exclude Patterns",
		ActionMessages: []string{"Checking Action Exclude Patterns"},
		StartedAt:      time.Now(),
		ParentID:       actionStepID,
		IsHidden:       true,
		Icon:           "solar:list-cross-minimalistic-bold",
	})
	if err != nil {
		log.Error("Error sending step:", err)
		return false, err
	}

	// end if there are no patterns
	if action.MatchPatterns[0].Key == "" {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             checkActionExcludeDataStep.ID,
			ActionMessages: []string{"Check skipped. Exclude patterns are disabled."},
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error("Error updating step:", err)
			return false, err
		}
		return false, nil
	}

	// convert payload to string
	payloadBytes, err := json.Marshal(payload.Payload)
	if err != nil {
		log.Error("Error converting payload to JSON:", err)
		return false, err
	}
	payloadString := string(payloadBytes)

	for _, excludePattern := range action.ExcludePatterns {
		value := gjson.Get(payloadString, excludePattern.Key)

		if value.String() == excludePattern.Value {
			err := executions.UpdateStep(execution, models.ExecutionSteps{
				ID:             checkActionExcludeDataStep.ID,
				ActionMessages: []string{`Match Pattern "` + excludePattern.Key + `": "` + excludePattern.Value + `" found`},
				NoPatternMatch: true,
				Finished:       true,
				FinishedAt:     time.Now(),
			})
			if err != nil {
				log.Error("Error updating step:", err)
				return false, err
			}
			return true, nil
		} else {
			err = executions.UpdateStep(execution, models.ExecutionSteps{
				ID:             checkActionExcludeDataStep.ID,
				ActionMessages: []string{`Exclude Pattern "` + excludePattern.Key + `": "` + excludePattern.Value + `" not found`},
			})
			if err != nil {
				log.Error("Error updating step:", err)
				return false, err
			}
		}
	}

	err = executions.UpdateStep(execution, models.ExecutionSteps{
		ID:             checkActionExcludeDataStep.ID,
		ActionMessages: []string{"Exclude Patterns checked"},
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step:", err)
		return false, err
	}

	return false, nil
}
