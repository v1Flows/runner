package flow

import (
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func CheckPatterns(flow models.Flows, execution models.Execution, payload models.Payload) (bool, error) {
	checkPatternsStep, err := executions.SendStep(execution, models.ExecutionSteps{
		ExecutionID:    execution.ID.String(),
		ActionName:     "Check Patterns",
		ActionMessages: []string{"Checking Patterns"},
		StartedAt:      time.Now(),
		Icon:           "solar:list-check-minimalistic-bold",
	})
	if err != nil {
		log.Error("Error sending step:", err)
		return false, err
	}

	// end if there are no patterns
	if len(flow.Patterns) == 0 {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             checkPatternsStep.ID,
			ActionMessages: []string{"Check skipped. No patterns are defined."},
			NoResult:       true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error("Error updating step:", err)
			return false, err
		}
		return true, nil
	}

	// convert payload to string
	payloadBytes, err := json.Marshal(payload.Payload)
	if err != nil {
		log.Error("Error converting payload to JSON:", err)
		return false, err
	}
	payloadString := string(payloadBytes)

	patternMissMatched := 0

	for _, pattern := range flow.Patterns {
		value := gjson.Get(payloadString, pattern.Key)

		if pattern.Type == "equals" {
			if value.String() == pattern.Value {
				err := executions.UpdateStep(execution, models.ExecutionSteps{
					ID:             checkPatternsStep.ID,
					ActionMessages: []string{`Pattern: ` + pattern.Key + ` == ` + pattern.Value + ` matched. Continue to next step`},
				})
				if err != nil {
					log.Error("Error updating step:", err)
					return false, err
				}
			} else {
				err = executions.UpdateStep(execution, models.ExecutionSteps{
					ID:             checkPatternsStep.ID,
					ActionMessages: []string{`Pattern: ` + pattern.Key + ` == ` + pattern.Value + ` not found. Skipping execution`},
					NoPatternMatch: true,
					Finished:       true,
					FinishedAt:     time.Now(),
				})
				if err != nil {
					log.Error("Error updating step:", err)
					return false, err
				}
				patternMissMatched++
			}
		} else if pattern.Type == "not_equals" {
			if value.String() != pattern.Value {
				err := executions.UpdateStep(execution, models.ExecutionSteps{
					ID:             checkPatternsStep.ID,
					ActionMessages: []string{`Pattern: ` + pattern.Key + ` != ` + pattern.Value + ` not found. Continue to next step`},
				})
				if err != nil {
					log.Error("Error updating step:", err)
					return false, err
				}
			} else {
				err = executions.UpdateStep(execution, models.ExecutionSteps{
					ID:             checkPatternsStep.ID,
					ActionMessages: []string{`Pattern: ` + pattern.Key + ` != ` + pattern.Value + ` matched. Skipping execution`},
					NoPatternMatch: true,
					Finished:       true,
					FinishedAt:     time.Now(),
				})
				if err != nil {
					log.Error("Error updating step:", err)
					return false, err
				}
				patternMissMatched++
			}
		}
	}

	if patternMissMatched > 0 {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             checkPatternsStep.ID,
			ActionMessages: []string{"Some patterns did not match. Skipping execution"},
			NoPatternMatch: true,
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error("Error updating step:", err)
			return false, err
		}
		executions.EndWithNoMatch(execution)
		return false, nil
	} else {
		err = executions.UpdateStep(execution, models.ExecutionSteps{
			ID:             checkPatternsStep.ID,
			ActionMessages: []string{"All patterns matched. Continue to next step"},
			Finished:       true,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error("Error updating step:", err)
			return false, err
		}
		return true, nil
	}
}
